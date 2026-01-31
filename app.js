var gameInterval
var player

document.addEventListener('DOMContentLoaded', function() {
   player = parseInt(new URLSearchParams(window.location.search).get('player'));
   render();
   const ws = new WebSocket("ws://localhost:9010/ws");

    ws.onmessage = (event) => {
       render()
    };

    ws.onclose = () => {
        console.log("connection closed");
    };
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
        renderPlayer(data);
        renderTable(data);
        renderStatus(data);
    })
    .catch(error => {
        console.error('Error:', error);
        clearInterval(gameInterval);
    });

    
}

function renderStatus(data) {
    const playerInfoElement = document.getElementById('playerInfo');
    const currentPlayer = getPlayer(player, data.Players);
    playerInfoElement.innerHTML = `${currentPlayer.Name}`;

    

    const statusElement = document.getElementById('status');
    statusElement.classList.remove('your-turn');
    statusElement.innerHTML = '';
    if(data.NextPlayer === player) {
        statusElement.innerHTML = `Your Turn!`;
        statusElement.classList.add('your-turn');
    }

    document.getElementById('getTrick').style.display = 'none';

    if(data.TrickWinner !== -1 && player === data.TrickWinner){
        document.getElementById('getTrick').style.display = 'block';
    }
    

}

function renderTable(data){
    const tableElement = document.getElementById('table');
    tableElement.innerHTML = '';
    if(data.Table === null) {
        return;
    }
    for (const card of data.Table) {
        li = document.createElement('li');
        li.textContent = `${card.Suit} ${card.Rank}`;
        tableElement.appendChild(li);
    }

    if(data.Table !== null && data.Table.length === 4){
        document.getElementById('getTrick').removeAttribute('disabled');
    } else {
        document.getElementById('getTrick').setAttribute('disabled', 'true');   
    }
}

function renderPlayer(data){
    if( data.Hand === null ){
        return;
    }
    const handElement = document.getElementById('hand');
    handElement.innerHTML = '';
    for (const card of data.Hand) {
        li = document.createElement('li');
        li.textContent = `${card.Suit} ${card.Rank}`;
        li.id = `${card.Id}`;
        if(card.Playable){
            li.classList.add('playable');
            li.addEventListener('click', playCard);
        }
        handElement.appendChild(li);
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

document.getElementById('getTrick').addEventListener('click', getTrick);

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