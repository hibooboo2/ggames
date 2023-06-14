function makeID(length) {
    let result = ''
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    const charactersLength = characters.length
    let counter = 0
    while (counter < length) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength))
        counter += 1
    }
    return result
}

function walkClickableElements() {
    // Get all elements on the page
    const allElements = document.querySelectorAll('*')

    // Array to store clickable elements
    const clickableElements = []

    // Iterate through each element
    allElements.forEach(element => {
        // Check if element is clickable

        const isClickable =
            element.tagName === 'A' || // Anchor tag
            element.tagName === 'BUTTON' || // Button tag
            element.tagName === 'INPUT' || // Input tag
            element.tagName === 'SELECT' || // Select tag
            element.tagName === 'TEXTAREA' || element.onclick != null // Textarea tag

        if (isClickable) {
            if (element.id == "" || element.id == null) {
                element.id = makeID(36)
            }
            clickableElements.push(element)
        }
    })

    return clickableElements
}

function clickElementById(id) {
    const element = document.getElementById(id)
    if (element) {
        element.click()
    } else {
        console.log(`Element with ID "${id}" not found.`)
    }
}

function getElementPosition(elementId) {
    const element = document.getElementById(elementId)

    if (element) {
        const position = element.getBoundingClientRect()
        return {
            top: position.top + window.scrollY,
            left: position.left + window.scrollX,
            bottom: position.bottom + window.scrollY,
            right: position.right + window.scrollX,
            width: position.width,
            height: position.height
        }
    } else {
        console.log(`Element with ID "${elementId}" not found.`)
        return null
    }
}

var shortcutDivs = []

function createHotkeyMarkers() {
    // Get all clickable elements on the page
    const clickableElements = walkClickableElements()

    // Create a new <div> element
    const div = document.createElement('div')
    div.className = 'hotkeyMarker'
    div.style.backgroundColor = 'yellow'
    div.style.position = 'fixed'
    div.style.width = '1.5em'
    div.style.height = '1em' // Adjust the height as needed
    div.style.zIndex = '9999'
    div.text = 'Click me!'

    // Insert the <div> element just above each clickable element
    clickableElements.forEach(element => {
        const parent = element.parentNode
        divForShortcut = div.cloneNode(true)
        pos = getElementPosition(element.id)
        console.log(pos)
        switch (element.className) {
            case "handToken":
                divForShortcut.style.top = (pos.top - 5) + 'px'
                divForShortcut.style.left = (pos.left - 5) + 'px'
                divForShortcut.style.width = (pos.width / 3) + 'px'
                divForShortcut.style.height = '12px'
                break
            case "playableCard":
                divForShortcut.style.top = (pos.top + 15) + 'px'
                divForShortcut.style.left = (pos.left + 15) + 'px'
                divForShortcut.style.width = (pos.width / 3) + 'px'
                divForShortcut.style.height = '12px'
                break
            default:
                divForShortcut.style.top = (pos.top - 5) + 'px'
                divForShortcut.style.left = (pos.left - 5) + 'px'
                break
        }
        shortcutDivs.push(divForShortcut)

        id = numberToUniqueLetters(shortcutDivs.length)
        shortcutToDivID[id] = element.id
        divForShortcut.textContent = id
        parent.insertBefore(divForShortcut, element)
    })
    shouldDetectHotKey = true
}

var existingIds = []

function numberToUniqueLetters(number) {
    const baseCharCode = 'A'.charCodeAt(0);
    const alphabetLength = 26;

    let result = '';

    while (number > 0) {
        number--; // Decrement by 1 to match 0-based index

        const remainder = number % alphabetLength;
        const letter = String.fromCharCode(baseCharCode + remainder);

        if (!existingIds.includes(letter)) {
            result = letter + result;
        }

        number = Math.floor(number / alphabetLength);
    }

    existingIds.push(result)
    return result;
}


var shortcutToDivID = {}
var sequencePressed = ""
var shouldDetectHotKey = false

function removeHotKeyMarkers() {
    console.log('Removing all hotkey boxes')
    shortcutDivs.forEach(function (el) {
        el.remove();
    });
    shortcutDivs = []
    shouldDetectHotKey = false
    existingIds = []
}

function hotKeyToggle(event) {
    console.log(event.key)
    if (event.key === "`") {
        if (shortcutDivs.length > 0) {
            removeHotKeyMarkers()
        } else {
            createHotkeyMarkers()
        }
        return
    }
    if (!shouldDetectHotKey) {
        return
    }
    sequencePressed += event.key
    sequencePressed = sequencePressed.toUpperCase()
    console.log(sequencePressed)
    if (sequencePressed in shortcutToDivID) {
        console.log(sequencePressed + " Clicking on element" + shortcutToDivID[sequencePressed])
        clickElementById(shortcutToDivID[sequencePressed])
        removeHotKeyMarkers()
        sequencePressed = ""
    }
}

document.addEventListener('keypress', hotKeyToggle)
