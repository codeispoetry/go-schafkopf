var player

document.addEventListener('DOMContentLoaded', function() {
   player = parseInt(new URLSearchParams(window.location.search).get('player'));

   if(isNaN(player)){
       document.body.innerHTML = 'Ungültiger Spieler';
       return;
   }
   render();
   const ws = new WebSocket("ws://localhost:9010/ws");

    ws.onmessage = (event) => {
       render()
    };

    ws.onclose = () => {
        console.log("connection closed");
    };

   
    document.querySelector('body').classList.add(`player${player}`);
});

/////////////////////////////
// Renderers
/////////////////////////////
function render(){
    fetch('http://localhost:9010/render', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "player": player
        })
    })
    .then(response => response.json())
    .then(data => {

        renderHand(data);
        renderTable(data);
        renderStatus(data);

        if(data.Status === "finished"){
            renderFinished(data);
        }
    })
    .catch(error => {
        console.error('Error:', error);
        document.body.innerHTML = 'Spielerver nicht erreichbar';
    });

    
}

function renderFinished(data){
    const status = document.querySelector('#status');
    status.innerHTML = '';

    status.innerHTML += `<br>${data.FinishLine}`;
}

function renderStatus(data) {
    const body = document.querySelector('body');

    if(data.Status !== "finished") {
        const status = document.querySelector('#status');
        status.innerHTML = '';
    }

    body.classList.remove('your-turn');
    if(data.NextPlayer === player) {
        body.classList.add('your-turn');
    }
}

function renderTable(data){
    const trickElement = document.getElementById('trick');
    const feltElement = document.getElementById('felt');

    feltElement.className = '';

    trickElement.innerHTML = '';
    if(data.Table === null) {
        return;
    }

    if(data.TrickWinner >= 0){
        feltElement.classList.add(`player${data.TrickWinner}`);
    }

    for (const card of data.Table) {
        li = document.createElement('li');
        li.setAttribute('title', `${card.Suit} ${card.Rank}`);
        li.classList.add(`${card.Suit.toLowerCase()}`);
        li.classList.add(`${card.Rank.toLowerCase().replace('ö', 'oe')}`);
        li.classList.add(`player${card.Player}`);
        li.style.zIndex = `${10 + card.Position}`;
        trickElement.appendChild(li);
    }

    document.getElementById('trick').dataset.takable = 'false';
    if(data.TrickWinner !== -1 && player === data.TrickWinner){
        document.getElementById('trick').dataset.takable = 'true';
    }

}

function renderHand(data){
    const handElement = document.getElementById('hand');
    handElement.innerHTML = '';
    
    if( data.Hand === null ){
        return;
    }

    let i = 0;
    for (const card of data.Hand) {
        li = document.createElement('li');
        li.setAttribute('title', `${card.Suit} ${card.Rank}`);
        li.id = `${card.Id}`;
        li.classList.add(`${card.Suit.toLowerCase()}`);
        li.classList.add(`${card.Rank.toLowerCase().replace('ö', 'oe')}`);

        const rotation = i * 5 - (data.Hand.length -1 ) * 5 /2;
        li.style.transform = `translateY(${Math.abs(rotation)*2}px)`;
        li.style.rotate = `${rotation}deg`;


        if(card.Playable){
            li.classList.add('playable');
            li.addEventListener('click', playCard);
        }
        handElement.appendChild(li);
        i++;
    }
}

//////////////////////////////
// Helpers
//////////////////////////////
function getPlayer(id, players){
    for(const player of players){
        if(player.Id === id){
            return {Name: player.Name};
        }
    }
    return {Name: "Unknown"};
}


///////////////////////////////
// Actions
///////////////////////////////
function playCard(event) {
    fetch('http://localhost:9010/play', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "player": player,
            "card": parseInt(event.target.id)
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        render()
       
    })
    .catch(error => console.error('Error:', error));
}

document.getElementById('trick').addEventListener('click', getTrick);
function getTrick(){
     fetch('http://localhost:9010/trick', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "player": player
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        render()
       
    })
    .catch(error => console.error('Error:', error));
}

document.addEventListener('click', (event) => {
    if (event.target.matches('#shuffle')) {
        startGame();
    }
});
function startGame(){
     fetch('http://localhost:9010/start', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({player: player})
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        render()
       
    })
    .catch(error => console.error('Error:', error));
}

document.querySelectorAll('[data-game]').forEach(button => {
    button.addEventListener('click', defineGame);
})
function defineGame(event){
    fetch('http://localhost:9010/define', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },          
        body: JSON.stringify({
            "player": player,
            "game": event.target.getAttribute('data-game'),
            "suit": event.target.getAttribute('data-suit')
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        render()
       
    })
    .catch(error => console.error('Error:', error));
}

function showButtons(data){

}