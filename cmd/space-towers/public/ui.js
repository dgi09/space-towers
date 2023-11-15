class UI {
    constructor(){
        this.visibleUI = "";

        let self = this;
        document.getElementById("btnStartRound").addEventListener("click", () => {
            self.onStartRoundHandler();
        });
    }

    onStartRound(callback){
        this.onStartRoundHandler = callback;
    }

    hideCurrentVisible(){
        if(this.visibleUI != "") {
            document.getElementById(this.visibleUI).className = "overlay hidden";
        }

        this.visibleUI = "";
       
        document.getElementById("canvas-overlay").className = "canvas-overlay hidden";
    }

    hideAll(){
        this.hideCurrentVisible();
    }

    showAwaitingRoundStart(scoreArray){
        this.hideCurrentVisible();
        this.visibleUI = "overlay-awaiting-round-start";

        document.getElementById("olsTitle").innerHTML = "Score";
        document.getElementById("btnStartRound").className = "";

        for(let i = 0; i < scoreArray.length; i++){
            let score = scoreArray[i];
            let playerNameElement = document.getElementById("olPlayer" + (i + 1) + "Name");
            playerNameElement.innerHTML = score.Player;

            let playerScoreElement = document.getElementById("olPlayer" + (i + 1) + "Score");
            playerScoreElement.innerHTML = score.Score;
        }

        
        document.getElementById("canvas-overlay").className = "canvas-overlay";
        document.getElementById(this.visibleUI).className = "overlay overlay-score";
    }

    showGameEnd(scoreArray){
        this.hideCurrentVisible();
        this.visibleUI = "overlay-awaiting-round-start";

        document.getElementById("olsTitle").innerHTML = "Game Ended";
        document.getElementById("btnStartRound").className = "hidden";
        
        for(let i = 0; i < scoreArray.length; i++){
            let score = scoreArray[i];
            let playerNameElement = document.getElementById("olPlayer" + (i + 1) + "Name");
            playerNameElement.innerHTML = score.Player;

            let playerScoreElement = document.getElementById("olPlayer" + (i + 1) + "Score");
            playerScoreElement.innerHTML = score.Score;
        }

        
        document.getElementById("canvas-overlay").className = "canvas-overlay";
        document.getElementById(this.visibleUI).className = "overlay overlay-score";
    }

    showRoundFinished(){
        this.hideCurrentVisible();
        this.visibleUI = "overlay-round-finished";

        document.getElementById("canvas-overlay").className = "canvas-overlay";
        document.getElementById(this.visibleUI).className = "overlay overlay-text";
    }

    showWaitingForOtherPlayers(){
        this.hideCurrentVisible();
        this.visibleUI = "overlay-waiting-for-other-players";

        document.getElementById("canvas-overlay").className = "canvas-overlay";
        document.getElementById(this.visibleUI).className = "overlay overlay-text";
    }
}