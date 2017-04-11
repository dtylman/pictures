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
)

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
	if !u.options.WithLocation {
		return nil
	}
	if image.HasPhase(picture.PhaseLocation) {
		return nil
	}
	err := location.PopulateLocation(image)
	if err != nil {
		tasklog.StatusMessage(tasklog.IndexerTask, err.Error())
	} else {
		tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Found location  %s", image.Location))
	}
	image.SetPhase(picture.PhaseLocation)
	return nil
}

func (u*Updater) thumbNail(image *picture.Index) error {
	if image.HasPhase(picture.PhaseThumb) {
		return nil
	}
	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Thumbing %s", image.Path))
	_, err := thumbs.MakeThumb(image.Path, image.MD5, u.options.DeleteDatabase)
	if err != nil {
		tasklog.StatusMessage(tasklog.IndexerTask, err.Error())
	}
	image.SetPhase(picture.PhaseThumb)
	return nil
}

func (u*Updater) objects(dp*darknet.Process, image*picture.Index) error {
	if !u.options.WithObjects {
		return nil
	}
	if image.HasPhase(picture.PhaseObjects) {
		return nil
	}
	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Detecing objects for %s", image.Path))
	res, err := dp.Detect(image.Path)
	image.SetPhase(picture.PhaseObjects)
	if err != nil {
		return err
	}
	if res.Result != darknet.Success {
		return errors.New(res.Result)
	}
	image.Objects = res.Objects
	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Found %v", image.Objects))

	return nil
}

func (u*Updater) IsEmpty() bool {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.images.Length() == 0
}

func (u*Updater) NextImage() *picture.Index {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	if u.images.Length() == 0 {
		return nil
	}
	return u.images.Remove().(*picture.Index)
}

func (u*Updater) worker(wg*sync.WaitGroup) {
	defer wg.Done()
	var dp *darknet.Process
	var err error
	if u.options.WithObjects {
		dp, err = darknet.NewProcess()
		if err != nil {
			tasklog.StatusMessage(tasklog.IndexerTask, err.Error())
		}
	}
	defer dp.Close()
	i := u.NextImage()
	for (i != nil) {
		err = u.thumbNail(i)
		if err != nil {
			tasklog.Error(err)
		}
		err = u.location(i)
		if err != nil {
			tasklog.Error(err)
		}
		err = u.objects(dp, i)
		if err != nil {
			tasklog.Error(err)
		}
		err = db.Index(i)
		if err != nil {
			tasklog.Error(err)
		}
		i = u.NextImage()
	}
}

func (u*Updater) update() {
	for !u.IsEmpty() {
		waitGroup := new(sync.WaitGroup)
		for i := 0; i < runtime.NumCPU(); i++ {
			waitGroup.Add(1)
			go u.worker(waitGroup)
		}
		waitGroup.Wait()
	}
}