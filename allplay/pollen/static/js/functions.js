function newGame() {
    let req = new XMLHttpRequest();
    req.open("POST", "/game");
    req.send();

    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        document.getElementById("mainbox").innerHTML = req.responseText
    }
}

function startGame(gameID) {
    let req = new XMLHttpRequest()
    req.open("POST", "/game/" + gameID + "/start/")
    req.send()

    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        document.getElementById("mainbox").innerHTML = req.responseText
        window.location.href = "/game/" + gameID + "/play/"
    }
}

function createJoinGameLink(gameID) {
    let req = new XMLHttpRequest()
    req.open("POST", "/tempID/?uname=" + document.getElementById("username").value + "&psw=" + document.getElementById("password").value)
    req.send()
    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        joinGameLink = location.protocol + '//' + location.host + "/game/" + gameID + "/join?anonymous=" + req.responseText
        navigator.clipboard.writeText(joinGameLink)
    }
}

function inviteUser(gameID) {
    let req = new XMLHttpRequest()
    req.open("POST", "/game/" + gameID + "/invite?username=" + document.getElementById("username").value)
    req.send()
    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        document.getElementById("username").value = ""
    }
}

function renderGame(gameID) {
    const source = new EventSource("/game/" + gameID + "/render/")
    source.onmessage = (event) => {
        switch (event.data) {
            case "game_not_started":
                document.getElementById("gamebox").innerHTML = "<h2>Waiting on game to start and render...</h2>"
            case "waiting":
                document.getElementById("gamebox").innerHTML = "<h2>Waiting on game to render...</h2>"
                return
        }
        console.log("Rendering board")
        cardToPlay = ""
        document.getElementById("gamebox").innerHTML = atob(event.data)

        document.removeEventListener('keypress', hotKeyToggle)
        document.addEventListener('keypress', hotKeyToggle)
    }
}

function playCard(gameID, cardID, x, y) {
    if (cardID == "") {
        alert("You must select a card to play")
        return
    }
    console.log("Play Card Called:", gameID, x, y)
    let req = new XMLHttpRequest()
    req.open("POST", "/game/" + gameID + "/play/card/" + cardID + "?position=" + x + ":" + y)
    req.send()
    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        console.log(req.responseText)
        cardToPlay = ""
    }
    //         /game/{ { $gameid } } /play/card / {{ $gameid }}
    // ?position = {{ $position.X }}: { { $position.Y } }
}

var cardToPlay = ""

function getCardToPlay() {
    return cardToPlay
}

function setCardToPlay(cardID) {
    if (cardToPlay != "") {
        document.getElementById(cardToPlay).className = "handCard"
        document.getElementById(cardToPlay + "_img").className = "handCard"
    }
    cardToPlay = cardID
    document.getElementById(cardToPlay).className = "handCardSelected"
    document.getElementById(cardToPlay + "_img").className = "handCardSelected"

}

var tokenToPlay = ""

function setTokenToPlay(tokenID) {
    if (tokenToPlay != "") {
        document.getElementById(tokenToPlay).className = "handToken"
        document.getElementById(tokenToPlay + "_img").className = "handToken"
    }
    tokenToPlay = tokenID
    document.getElementById(tokenToPlay).className = "handTokenSelected"
    document.getElementById(tokenToPlay + "_img").className = "handTokenSelected"
}

function playToken(gameID, x, y) {
    console.log("Play Token Called:", gameID, x, y)
    if (tokenToPlay == "") {
        alert("You must select a token to play")
        return
    }
    console.log("Play Token Called:", gameID, x, y)
    let req = new XMLHttpRequest()
    req.open("POST", "/game/" + gameID + "/play/token/" + tokenToPlay + "?position=" + x + ":" + y)
    req.send()
    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        console.log(req.responseText)
        tokenToPlay = ""
    }
    //         /game/{ { $gameid } } /play/card / {{ $gameid }}
    // ?position = {{ $position.X }}: { { $position.Y } }
}

function toggleHints(gameID) {
    let req = new XMLHttpRequest()
    req.open("POST", "/game/" + gameID + "/hints/toggle/")
    req.send()
}
