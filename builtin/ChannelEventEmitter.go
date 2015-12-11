package builtin

import (
	"github.com/robertkrimen/otto"
	"github.com/vijayee/Contract"
)

var ChannelEventEmitter contract.API
var EmitterScript = `
	var EmitterRegistry=[];
	var EmitterId= (function(){
		var current;
		function init(){
			return 0;
		}
		return{
			next: function(){
				if(!current){
					current= inti();
				}
				return ++current;
			}
		};
	})();	
	
	var NewEmitter= function(obj){
		obj.prototype._EmitterId = EmitterId.next();
		obj.prototype.emit= function(evt, data){
			Emit(this._EmitterId, evt, data);
		};
		EmitterRegistry.push(obj);
		return obj;
	};

`

type Event struct {
	name  string
	value otto.Value
	from  int64
}

var subscribers map[int64]map[string]map[int64]chan Event // ehrmahgaud this is long
var index int64

func Emit(call otto.FunctionCall, conv contract.Converter) otto.Value {
	if len(call.ArgumentList) < 3 {
		return otto.Value{}
	}
	jsObjectId := call.ArgumentList[0]
	jsEventName := call.ArgumentList[1]
	jsValue := call.ArgumentList[2]

	objectid, err := jsObjectId.ToInteger()
	if err != nil {
		return otto.Value{}
	}
	name, err := jsEventName.ToString()
	if err != nil {
		return otto.Value{}
	}
	EmitOnChannels(Event{name: name, from: objectid, value: jsValue})
	return otto.Value{}

}

func EmitOnChannels(event Event) {
	object, ok := subscribers[event.from]
	if ok == false {
		return
	}
	events, ok := object[event.name]
	if ok == false {
		return
	}
	for _, listener := range events {
		select {
		case listener <- event:
			continue
		default:
			continue
		}
	}
}
func Subscribe(objectid int64, name string, channel chan Event) {
	if subscribers == nil {
		subscribers = make(map[int64]map[string]map[int64]chan Event)
	}
	object, ok := subscribers[objectid]
	if ok == false {
		object = make(map[string]map[int64]chan Event)
	}
	event, ok := object[name]
	if ok == false {
		event = make(map[int64]chan Event)
		event[index] = channel
		index++
	} else {
		found := false
		for _, sub := range event {
			if sub == channel {
				found = true
				break
			}
		}
		if found == false {
			event[index] = channel
			index++
		}
	}
	object[name] = event
	subscribers[objectid] = object
}

func UnSubscribe(objectid int64, name string, channel chan Event) {
	if subscribers == nil {
		return
	} else {
		object, ok := subscribers[objectid]
		if ok == false {
			return
		} else {
			event, ok := object[name]
			if ok == false {
				return
			}
			var key int64
			found := false
			for i, listener := range event {
				if listener == channel {
					found = true
					key = i
				}
			}
			if found == true {
				delete(event, key)
			}
			if len(event) == 0 {
				delete(object, name)
			} else {
				object[name] = event
			}
			if len(object) == 0 {
				delete(subscribers, objectid)
			} else {
				subscribers[objectid] = object
			}
		}
	}
}
