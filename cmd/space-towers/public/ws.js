const IN_TYPE_READY = 0;
const IN_TYPE_PREPARE_GAME = 1;
const IN_TYPE_AWAIT_ROUND_START = 2;
const IN_TYPE_START_ROUND = 3;
const IN_TYPE_GAME_ENDED = 4;
const IN_TYPE_FORCE_ROUND_FINISH = 5;

const OUT_TYPE_DETAILS = 0;
const OUT_TYPE_CREATE_ROOM = 1;
const OUT_TYPE_JOIN_ROOM = 2;
const OUT_TYPE_GAME_READY = 3;
const OUT_TYPE_START_ROUND = 4;
const OUT_TYPE_ROUND_FINISHED = 5;

class Ws {
  constructor(url) {
    this.url = url;
    this.ws = new WebSocket(url);

    let self = this;
    this.ws.onmessage = (event) => {
      self.handleMsg(event.data);
    };
  }

  onOpen(callback) {
    this.ws.onopen = callback;
  }

  onReady(callback) {
    this.onReadyMsg = callback;
  }

  onPrepareGame(callback) {
    this.onPrepareGameMsg = callback;
  }

  onAwaitRoundStart(callback) {
    this.onAwaitRoundStartMsg = callback;
  }

  onStartRound(callback){
    this.onStartRoundMsg = callback;
  }

  onGameEnded(callback){
    this.onGameEndedMsg = callback;
  }

  onForceRoundFinish(callback){
    this.onForceRoundFinishMsg = callback;
  }

  handleMsg(data) {
    let self = this;
    let buf = data.arrayBuffer().then(b => {
      let array = new Int8Array(b);
      let type = array[0];

      let msg = null;
      if (array.length > 1) {
        msg = JSON.parse(String.fromCharCode.apply(null, array.slice(1)));
      }

      switch (type) {
        case IN_TYPE_READY:
          self.onReadyMsg();
          break;

        case IN_TYPE_PREPARE_GAME:
          self.onPrepareGameMsg(msg);
          break;

        case IN_TYPE_AWAIT_ROUND_START:
          self.onAwaitRoundStartMsg(msg);
          break;

        case IN_TYPE_START_ROUND:
          self.onStartRoundMsg(msg);
          break;

        case IN_TYPE_GAME_ENDED:
          self.onGameEndedMsg(msg);
          break;

        case IN_TYPE_FORCE_ROUND_FINISH:
          self.onForceRoundFinishMsg(msg);
          break;
      }
    });
  }

  send(type, msg) {
    let msgLen = 0;
    let msgStr = null;

    if (msg != null) {
      msgStr = JSON.stringify(msg);
      msgLen = msgStr.length;
    }

    let buffer = new ArrayBuffer(1 + msgLen);
    let array = new Int8Array(buffer);

    array[0] = type;

    if (msgLen > 0) {
      for (let i = 0; i < msgStr.length; i++) {
        array[i + 1] = msgStr.charCodeAt(i);
      }
    }

    this.ws.send(buffer);
  }

  sendDetails(name) {
    this.send(OUT_TYPE_DETAILS, { name: name });
  }

  sendCreateRoom(room) {
    this.send(OUT_TYPE_CREATE_ROOM, { room: room });
  }

  sendJoinRoom(room){
    this.send(OUT_TYPE_JOIN_ROOM, { room: room });
  }

  sendQueueForGame() {
    this.send(OUT_TYPE_QUEUE_FOR_GAME, null);
  }

  sendGameReady(){
    this.send(OUT_TYPE_GAME_READY, null);
  }

  sendStartRound(){
    this.send(OUT_TYPE_START_ROUND, null);
  }

  sendRoundFinished(combos){
    this.send(OUT_TYPE_ROUND_FINISHED, { combos: combos });
  }
}