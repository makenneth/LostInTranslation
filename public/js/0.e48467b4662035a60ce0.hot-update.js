webpackHotUpdate(0,{

/***/ 317:
/***/ function(module, exports) {

	"use strict";
	
	Object.defineProperty(exports, "__esModule", {
	  value: true
	});
	var TOGGLE_CHAT = "hangperson/chat-screen/TOGGLE";
	
	exports.default = function () {
	  var state = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : false;
	  var action = arguments[1];
	
	  if (action.type === TOGGLE_CHAT) {
	    return !state;
	  }
	
	  return state;
	};
	
	var toggleChat = exports.toggleChat = function toggleChat() {
	  return {
	    type: TOGGLE_CHAT
	  };
	};

/***/ }

})
//# sourceMappingURL=0.e48467b4662035a60ce0.hot-update.js.map