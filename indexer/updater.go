package indexer

import (
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/tasklog"
	"fmt"
	"github.com/dtylman/pictures/indexer/thumbs"
	"github.com/dtylman/pictures/indexer/location"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/eapache/queue"
	"runtime"
	"sync"
	"github.com/dtylman/pictures/indexer/darknet"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/dtylman/pictures/conf"
	"time"
)

var darknetError = errors.New("Darknet Process Error")

type Updater struct {
	images  *queue.Queue
	options Options
	mutex   sync.Mutex
}

func NewUpdater(options Options) *Updater {
	p := &Updater{
		images: queue.New(),
		options: options,
	}
	db.WalkImages(p.walkImage)
	return p
}

func (u*Updater) walkImage(key string, image *picture.Index, err error) {
	u.images.Add(image)
}

func (u*Updater) location(image*picture.Index) error {
	if image.HasPhase(picture.PhaseLocation) {
		return nil
	}
	defer image.SetPhase(picture.PhaseLocation)

	err := location.PopulateLocation(image)
	if err != nil {
		tasklog.Error(err)
	} else {
		tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Found location  %s", image.Location))
	}
	return nil
}

func (u*Updater) thumbNail(image *picture.Index) error {
	if image.HasPhase(picture.PhaseThumb) {
		return nil
	}
	defer image.SetPhase(picture.PhaseThumb)

	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Thumbing %s", image.Path))
	_, err := thumbs.MakeThumb(image.Path, image.MD5, u.options.DeleteDatabase)
	if err != nil {
		tasklog.StatusMessage(tasklog.IndexerTask, err.Error())
	}
	return nil
}

func (u*Updater) objects(dp*darknet.Process, image*picture.Index) error {
	if image.HasPhase(picture.PhaseObjects) {
		return nil
	}
	defer image.SetPhase(picture.PhaseObjects)

	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Detecing objects for %s", image.Path))
	res, err := dp.Detect(image.Path, time.Duration(conf.Options.DarknetTimeout) * time.Second)
	if err != nil {
		tasklog.Error(err)
		return darknetError
	}
	if res.Result != darknet.Success {
		return errors.New(res.Result)
	}
	for _, o := range res.Objects {
		image.Objects += fmt.Sprintf("%d %s with %d %% ", o.Count, o.Name, o.Prob)
	}
	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Found %v", res.Objects))
	return nil
}

func (u*Updater) Length() int {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.images.Length()
}

func (u*Updater) NextImage() (int, *picture.Index) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	size := u.images.Length()
	if size == 0 {
		return 0, nil
	}
	return size, u.images.Remove().(*picture.Index)
}

func (u*Updater) getDarknessProcess(dp*darknet.Process) *darknet.Process {
	if dp != nil {
		dp.Close()
	}
	if u.options.WithObjects {
		var err error
		dp, err = darknet.NewProcess()
		if err != nil {
			tasklog.Error(err)
		}
		return dp
	}
	return nil
}

func (u*Updater) worker(wg*sync.WaitGroup, total int) {
	defer wg.Done()
	var dp *darknet.Process
	var err error

	defer func() {
		if dp != nil {
			dp.Close()
		}
	}()
	left, i := u.NextImage()
	for (i != nil) {
		tasklog.Status(tasklog.IndexerTask, IsRunning(), total - left, total, fmt.Sprintf("Indexing %s", i.Path))
		if !IsRunning() {
			//indexer had stopped.
			return
		}
		err = u.thumbNail(i)
		if err != nil {
			tasklog.Error(err)
		}
		if u.options.WithLocation {
			err = u.location(i)
			if err != nil {
				tasklog.Error(err)
			}
		}
		if u.options.WithObjects {
			if dp == nil {
				dp = u.getDarknessProcess(dp)
			}
			err = u.objects(dp, i)
			if err != nil {
				tasklog.Error(err)
				if err == darknetError {
					//respawn
					dp = u.getDarknessProcess(dp)
				}
			}

		}
		if u.options.WithFaces {

		}
		err = db.Index(i)
		if err != nil {
			tasklog.Error(err)
		}
		left, i = u.NextImage()
	}
}

func (u*Updater) update() {
	total := u.Length()
	for u.Length() > 0 {
		waitGroup := new(sync.WaitGroup)
		for i := 0; i < runtime.NumCPU(); i++ {
			waitGroup.Add(1)
			go u.worker(waitGroup, total)
		}
		waitGroup.Wait()
	}
}
