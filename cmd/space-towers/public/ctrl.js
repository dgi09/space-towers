function onConnectedToServer() {
    document.getElementById("stageAttemptingConnect").className = "hidden";
    document.getElementById("stageEnterDetails").className = "";
}

function onReady() {
    document.getElementById("stageEnterDetails").className = "hidden";
    document.getElementById("stageControlls").className = "";
}

function onPrepareGame(data) {
    let players = data.players;
    let maxRounds = data.maxRounds;
    let roundDuration = data.roundDuration;

    document.getElementById("stageControlls").className = "hidden";
    document.getElementById("canvas").className = "";

    game = new Game("canvas", ws, players, maxRounds, roundDuration);
    game.Start();
}

function onEnterDetailsClicked() {
    let name = document.getElementById("inputName").value;

    if (name === undefined || name.length == 0) {
        return;
    }

    ws.sendDetails(name);
}

function onQueueForGameClicked() {
    ws.sendQueueForGame();
}

function onCreateRoomClicked() {
    let room = document.getElementById("inputCreateRoomName").value;

    if (room === undefined || room.length == 0) {
        return;
    }

    ws.sendCreateRoom(room);
}


function onJoinRoomClicked() {
    let room = document.getElementById("inputJoinRoomName").value;

    if (room === undefined || room.length == 0) {
        return;
    }

    ws.sendJoinRoom(room);
}