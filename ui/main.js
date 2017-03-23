
var child;
var fails = 0;

function body_message(msg){    
    document.body.innerHTML = '<h1>'+msg+'</h1>';        
}

function start_process() {
    body_message("Loading...");

    const spawn = require('child_process').spawn;
    child = spawn('/home/danny/src/go/src/github.com/dtylman/pictures/cmd/pc/pc',{maxBuffer:1024*500});    
    
    const readline = require('readline');
    const rl = readline.createInterface({
        input: child.stdout
    })

    rl.on('line', (data) => {        
        console.log(`Received: ${data}`);
        document.body.innerHTML = data;
    });

    child.stderr.on('data', (data) => {
        console.log(`stderr: ${data}`);
    });

    child.on('close', (code) => {                
        body_message(`process exited with code ${code}`);
        restart_process();    
    });

    child.on('error', (err) => {
        body_message('Failed to start child process.');
        restart_process();
    });
}

function restart_process(){
    setTimeout(function () {
            fails++;
            if (fails > 3) {
                close();
            } else {
                start_process();
            }
        }, 1000);
}

function element_as_object(elem){
    var obj = {
        value : elem.value,
        attributes : {}
    }
    for (var j=0;j<elem.attributes.length;j++){
        obj.attributes[elem.attributes[j].name]=elem.attributes[j].value;         
    }
    return obj;
}

function element_by_tag_as_array(tag){
    var items = [];
    var elems = document.getElementsByTagName(tag); 
    for (var i = 0; i < elems.length; i++) {
        items.push(element_as_object(elems[i]));
    }    
    return items;
}

function fire_event(sender){
    var msg = {
        sender : element_as_object(sender),
        inputs : element_by_tag_as_array("input")
    }    
    child.stdin.write(JSON.stringify(msg));
    console.log(JSON.stringify(msg));
}

start_process();

