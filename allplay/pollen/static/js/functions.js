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
        renderGame(gameID)
    }
}

function createJoinGameLink(gameID) {
    let req = new XMLHttpRequest()
    req.open("POST", "/tempID/")
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
            case "waiting":
                document.getElementById("gamebox").innerHTML = "<h2>Waiting on game to start and render...</h2>"
                return
        }
        console.log("Rendering board")
        document.getElementById("gamebox").innerHTML = atob(event.data)
    }
}

function playCard(gameID, cardID, x, y) {
    console.log("Play Card Called:", gameID, x, y)
    let req = new XMLHttpRequest()
    req.open("POST", "/game/" + gameID + "/play/card/" + cardID + "?position=" + x + ":" + y)
    req.send()
    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        console.log(req.responseText)
    }
    //         /game/{ { $gameid } } /play/card / {{ $gameid }
    // }?position = {{ $position.X }}: { { $position.Y } }
}
