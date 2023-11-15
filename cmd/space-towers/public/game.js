const CARDS_GRAPH = [
    [
        [0],
        [0, 1],
        [1, 2],
        [2, 3],
        [3, 4],
        [4, 5],
        [5, 6],
        [6, 7],
        [7, 8],
        [8]
    ],
    [
        [0],
        [0, 1],
        [1],
        [2],
        [2, 3],
        [3],
        [4],
        [4, 5],
        [5],
    ],
    [
        [0],
        [0],
        [1],
        [1],
        [2],
        [2]
    ]
]

class Game {
    constructor(canvasId, ws, playersNames, maxRounds, roundDuration) {
        this.ui = new UI();
        this.canvasId = canvasId;
        this.canvas = document.getElementById(canvasId);
        this.ctx = this.canvas.getContext("2d");
        this.ws = ws;
        this.players = playersNames;
        this.maxRounds = maxRounds;
        this.roundDuration = roundDuration;

        this.w = canvas.width;
        this.h = canvas.height;

        this.cardImage = document.getElementById("imgCards");
        this.cardWidth = 319.14;
        this.cardHeight = 449;
        this.cardAR = this.cardWidth / this.cardHeight;
        this.cardsInRow = 14;

        this.cardVWidth = this.w / 12;
        this.cardVHeight = this.cardVWidth / this.cardAR;

        this.currentRoundId = 0;

        this.ws.onAwaitRoundStart(this.onAwaitRoundStartReceived.bind(this));
        this.ws.onStartRound(this.onStartRoundReceived.bind(this));
        this.ws.onGameEnded(this.onGameEndedReceived.bind(this));
        this.ws.onForceRoundFinish(this.onForceRoundFinishReceived.bind(this));

        this.ui.onStartRound(this.onStartRoundClicked.bind(this));

        this.canvas.addEventListener("click", this.onCanvasClick.bind(this));
        this.calcExtraPositions();
    }

    Start() {
        this.clear();

        this.ws.sendGameReady();
    }

    onAwaitRoundStartReceived(msg) {
        this.ui.showAwaitingRoundStart(msg.score);
    }

    onStartRoundClicked() {
        this.ui.showWaitingForOtherPlayers();

        this.ws.sendStartRound();
    }

    onStartRoundReceived(msg) {
        this.ui.hideAll();

        this.startRound(msg);
    }

    onGameEndedReceived(msg) {
        this.ui.showGameEnd(msg.score);
        this.clear();
    }

    onForceRoundFinishReceived() {
        this.ui.showRoundFinished();

        this.ws.sendRoundFinished(this.round.combos);
        this.clear();
    }

    onCanvasClick(e) {
        let rect = this.canvas.getBoundingClientRect();
        let x = e.clientX - rect.left;
        let y = e.clientY - rect.top;

        let px = e.clientX / (rect.right - rect.left) * this.w;
        let py = e.clientY / (rect.bottom - rect.top) * this.h;

        this.handleInteract(px, py);

        if (this.evalRoundFinished()) {
            this.ui.showRoundFinished();

            this.ws.sendRoundFinished(this.round.combos);
            this.clear();
            return;
        }

        this.draw();
    }

    handleInteract(x, y) {
        let pos = this.pos;

        if (this.clickInCardRect(x, y, pos.deck.x, pos.deck.y)) {
            if (this.round.extraIndex < this.round.extraCards.length) {
                this.round.extraIndex++;
                this.round.extraCard = this.round.extraCards[this.round.extraIndex];

                this.round.wildCard = null;
                this.round.combos.push(this.round.combo);
                this.round.combo = 0;
            }

            return;
        }

        for (let i = 0; i < this.round.row1.length; i++) {
            let entry = this.round.row1[i];

            if (this.clickInCardRect(x, y, entry.x, entry.y)) {
                if (entry.played) {
                    continue;
                }

                if (!entry.visible) {
                    return;
                }

                if (this.tryPlaceCard(entry.card)) {
                    entry.played = true;

                    for (let j = 0; j < entry.cardConnections.length; j++) {
                        let index = entry.cardConnections[j];
                        let otherEntry = this.round.row2[index];
                        otherEntry.cardBelow--;
                        if (otherEntry.cardBelow == 0) {
                            otherEntry.visible = true;
                        }
                    }
                }

                return;
            }
        }

        for (let i = 0; i < this.round.row2.length; i++) {
            let entry = this.round.row2[i];

            if (this.clickInCardRect(x, y, entry.x, entry.y)) {
                if (entry.played) {
                    continue;
                }

                if (!entry.visible) {
                    return;
                }

                if (this.tryPlaceCard(entry.card)) {
                    entry.played = true;

                    for (let j = 0; j < entry.cardConnections.length; j++) {
                        let index = entry.cardConnections[j];
                        let otherEntry = this.round.row3[index];
                        otherEntry.cardBelow--;
                        if (otherEntry.cardBelow == 0) {
                            otherEntry.visible = true;
                        }
                    }
                }

                return;
            }
        }

        for (let i = 0; i < this.round.row3.length; i++) {
            let entry = this.round.row3[i];

            if (this.clickInCardRect(x, y, entry.x, entry.y)) {
                if (entry.played) {
                    continue;
                }

                if (!entry.visible) {
                    return;
                }

                if (this.tryPlaceCard(entry.card)) {
                    entry.played = true;

                    for (let j = 0; j < entry.cardConnections.length; j++) {
                        let index = entry.cardConnections[j];
                        let otherEntry = this.round.row4[index];
                        otherEntry.cardBelow--;
                        if (otherEntry.cardBelow == 0) {
                            otherEntry.visible = true;
                        }
                    }
                }

                return;
            }
        }

        for (let i = 0; i < this.round.row4.length; i++) {
            let entry = this.round.row4[i];

            if (this.clickInCardRect(x, y, entry.x, entry.y)) {
                if (entry.played) {
                    continue;
                }

                if (!entry.visible) {
                    return;
                }

                if (this.tryPlaceCard(entry.card)) {
                    entry.played = true;
                }

                return;
            }
        }
    }

    tryPlaceCard(card) {
        if (this.round.wildCard != null) {
            if (this.areCardsCont(card, this.round.wildCard)) {
                this.round.wildCard = card;
                this.round.combo++;

                this.evalCombo();
                return true;
            }
        }

        if (this.areCardsCont(card, this.round.extraCard)) {
            this.round.extraCard = card;
            this.round.combo++;

            this.evalCombo();
            return true;
        }

        return false;
    }

    canPlaceCard(card) {
        if (this.round.wildCard != null) {
            if (this.areCardsCont(card, this.round.wildCard)) {
                return true;
            }
        }

        return this.areCardsCont(card, this.round.extraCard);
    }

    evalCombo() {
        if (this.round.combo >= 3 && this.round.wildCard == null) {
            if (this.round.extraIndex < this.round.extraCards.length) {
                this.round.extraIndex++;
                this.round.wildCard = this.round.extraCards[this.round.extraIndex];
            }
        }

    }

    evalRoundFinished() {
        let playedCards = 0;

        for (let i = 0; i < this.round.row1.length; i++) {
            let entry = this.round.row1[i];
            if (entry.played) {
                playedCards++;
            } else {
                if (entry.visible && this.canPlaceCard(entry.card)) {
                    return false;
                }
            }
        }

        for (let i = 0; i < this.round.row2.length; i++) {
            let entry = this.round.row2[i];
            if (entry.played) {
                playedCards++;
            } else {
                if (entry.visible && this.canPlaceCard(entry.card)) {
                    return false;
                }
            }
        }

        for (let i = 0; i < this.round.row3.length; i++) {
            let entry = this.round.row3[i];
            if (entry.played) {
                playedCards++;
            } else {
                if (entry.visible && this.canPlaceCard(entry.card)) {
                    return false;
                }
            }
        }

        for (let i = 0; i < this.round.row4.length; i++) {
            let entry = this.round.row4[i];
            if (entry.played) {
                playedCards++;
            } else {
                if (entry.visible && this.canPlaceCard(entry.card)) {
                    return false;
                }
            }
        }

        if (playedCards == 28) {
            return true;
        }

        return this.round.extraIndex == this.round.extraCards.length - 1;
    }

    areCardsCont(card1, card2) {
        let upRank = card1.Rank + 1;
        let downRank = card1.Rank - 1;

        if (upRank > 13) {
            upRank = 1;
        }

        if (downRank < 1) {
            downRank = 13;
        }

        return card2.Rank == upRank || card2.Rank == downRank;
    }

    clickInCardRect(x, y, cX, cY) {
        return x >= cX && x <= cX + this.cardVWidth && y >= cY && y <= cY + this.cardVHeight;
    }

    calcExtraPositions() {
        let pos = {};

        let gapBetDeckAndExtra = 10;
        let gapBetExtraAndWild = 1.7 * this.cardVWidth;

        let totalWidth = this.cardVWidth + gapBetDeckAndExtra + this.cardVWidth + gapBetExtraAndWild + this.cardVWidth;

        let startX = (this.w - totalWidth) / 2;
        let startY = this.h - this.cardVHeight - 5;

        pos.deck = { x: startX, y: startY };

        let extraCardsRemaining = 52;
        let textMeasurements = this.ctx.measureText(extraCardsRemaining.toString());
        pos.cardCount = { x: startX + this.cardVWidth / 2 - textMeasurements.width, y: startY + this.cardVHeight / 2 + 15 };

        startX += this.cardVWidth + gapBetDeckAndExtra;

        pos.extra = { x: startX, y: startY };

        startX += this.cardVWidth + gapBetExtraAndWild;

        pos.wild = { x: startX, y: startY };

        this.pos = pos;
    }

    startRound(data) {
        let round = {};

        round.id = data.round;
        round.duration = data.duration;

        let deck = data.deck;

        round.extraIndex = 0;
        round.extraCards = data.deck.Extra;
        round.extraCard = round.extraCards[round.extraIndex];
        round.wildCard = null;
        round.combo = 0;
        round.combos = [];

        round.row1 = [];
        for (let i = 0; i < deck.Row1.length; i++) {
            round.row1.push({
                cardBelow: 0,
                cardConnections: CARDS_GRAPH[0][i],
                played: false,
                visible: true,
                card: deck.Row1[i]
            });
        }

        round.row2 = [];
        for (let i = 0; i < deck.Row2.length; i++) {
            round.row2.push({
                cardBelow: 2,
                cardConnections: CARDS_GRAPH[1][i],
                played: false,
                visible: false,
                card: deck.Row2[i]
            });
        }

        round.row3 = [];
        for (let i = 0; i < deck.Row3.length; i++) {
            round.row3.push({
                cardBelow: 2,
                cardConnections: CARDS_GRAPH[2][i],
                played: false,
                visible: false,
                card: deck.Row3[i]
            });
        }

        round.row4 = [];
        for (let i = 0; i < deck.Row4.length; i++) {
            round.row4.push({
                cardBelow: 2,
                cardConnections: [],
                played: false,
                visible: false,
                card: deck.Row4[i]
            });
        }

        this.round = round;

        let halfCardW = this.cardVWidth / 2;
        let halfCardH = this.cardVHeight / 2;

        let totalW = this.cardVWidth * 10;
        let totalH = this.cardVHeight + halfCardH * 3;

        let startX = (this.w - totalW) / 2;
        let startY = (this.h - totalH) / 2 + this.cardVHeight;

        for (let i = 0; i < round.row1.length; i++) {
            let card = round.row1[i];
            card.x = startX + i * this.cardVWidth;
            card.y = startY;
        }

        startX += halfCardW;
        startY -= halfCardH;

        for (let i = 0; i < round.row2.length; i++) {
            let card = round.row2[i];
            card.x = startX + i * this.cardVWidth;
            card.y = startY;
        }

        startX += halfCardW;
        startY -= halfCardH;

        let gap = 0;
        for (let i = 0; i < round.row3.length; i++) {
            let card = round.row3[i];
            card.x = startX + (i + gap) * this.cardVWidth;
            card.y = startY;

            if ((i + 1) % 2 == 0) {
                gap++;
            }
        }

        startX += halfCardW;
        startY -= halfCardH;

        gap = 0;
        for (let i = 0; i < round.row4.length; i++) {
            let card = round.row4[i];
            card.x = startX + (i + gap) * this.cardVWidth;
            card.y = startY;

            gap += 2;
        }

        this.draw();
    }

    clear() {
        this.ctx.fillStyle = "#000000";
        this.ctx.fillRect(0, 0, this.w, this.h);
    }

    draw() {
        this.clear();
        this.drawBoard();
        this.drawExtraBoard();
    }

    drawBoard() {
        let round = this.round;
        for (let i = 0; i < round.row4.length; i++) {
            let entry = round.row4[i];
            if (!entry.played) {
                if (!entry.visible) {
                    this.drawCardBack(entry.x, entry.y);
                } else {
                    let card = entry.card;
                    this.drawCard(entry.x, entry.y, card);
                }
            }
        }

        for (let i = 0; i < round.row3.length; i++) {
            let entry = round.row3[i];
            if (!entry.played) {
                if (!entry.visible) {
                    this.drawCardBack(entry.x, entry.y);
                } else {
                    let card = entry.card;
                    this.drawCard(entry.x, entry.y, card);
                }
            }
        }

        for (let i = 0; i < round.row2.length; i++) {
            let entry = round.row2[i];
            if (!entry.played) {
                if (!entry.visible) {
                    this.drawCardBack(entry.x, entry.y);
                } else {
                    let card = entry.card;
                    this.drawCard(entry.x, entry.y, card);
                }
            }
        }

        for (let i = 0; i < round.row1.length; i++) {
            let entry = round.row1[i];
            if (!entry.played) {
                if (!entry.visible) {
                    this.drawCardBack(entry.x, entry.y);
                } else {
                    let card = entry.card;
                    this.drawCard(entry.x, entry.y, card);
                }
            }
        }
    }

    drawExtraBoard() {
        let round = this.round;
        let pos = this.pos;
        let extraCardsRemaining = round.extraCards.length - round.extraIndex - 1;
        if (extraCardsRemaining > 0) {
            this.drawCardBack(pos.deck.x, pos.deck.y);

            this.ctx.fillStyle = "#FFFFFF";
            this.ctx.font = "bold 30px Arial";
            this.ctx.fillText(extraCardsRemaining.toString(), pos.cardCount.x, pos.cardCount.y);
        }


        this.drawCard(pos.extra.x, pos.extra.y, round.extraCard);

        if (round.wildCard != null) {
            this.drawCard(pos.wild.x, pos.wild.y, round.wildCard);
        }
    }

    drawCard(x, y, card) {
        let suit = card.Suit;
        let rank = card.Rank;

        let row = suit;
        let col = rank - 1;

        this.ctx.drawImage(this.cardImage, col * this.cardWidth, row * this.cardHeight, this.cardWidth, this.cardHeight, x, y, this.cardVWidth, this.cardVHeight);
    }

    drawCardBack(x, y) {
        this.ctx.drawImage(this.cardImage, 13 * this.cardWidth, 1 * this.cardHeight, this.cardWidth, this.cardHeight, x, y, this.cardVWidth, this.cardVHeight);
    }
}