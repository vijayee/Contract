function(global){
	var go;
	if(global.go){
		go =global.go;
	}
	var emitter = new EventEmitter(); 
	util.inherits(go, emitter);
	
	
	
}.call(this,this)