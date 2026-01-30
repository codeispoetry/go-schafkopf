document.addEventListener('DOMContentLoaded', function() {
   setInterval(getGame, 1500);
   getGame();
});

function getGame(){
     fetch('http://localhost:9010/play', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "player": parseInt(new URLSearchParams(window.location.search).get('player'))
        })
    })
    .then(response => response.json())
    .then(data => {
        render(data);
    })
    .catch(error => console.error('Error:', error));
}

function render(data){
    renderPlayer(data);
    renderTable(data);
    renderStatus(data);
}

function renderStatus(data){
    const playerInfoElement = document.getElementById('playerInfo');
    const player = getPlayer(parseInt(new URLSearchParams(window.location.search).get('player')), data.Players);
    playerInfoElement.innerHTML = `You are: ${player.Name} with ${player.Points} points.`;

    const statusElement = document.getElementById('status');
    if(data.NextPlayer === parseInt(new URLSearchParams(window.location.search).get('player'))) {
        statusElement.innerHTML = `Your Turn!`;
        statusElement.classList.add('your-turn');
        return;
    }
    statusElement.classList.remove('your-turn');
    statusElement.innerHTML = `Next Player: ${getPlayer(data.NextPlayer, data.Players).Name}`;
}

function getPlayer(id, players){
    for(const player of players){
        if(player.Id === id){
            return {Name: player.Name, Points: player.Points};
        }
    }
    return {Name: "Unknown", Points: 0};
}

function renderPlayer(data){
    const handElement = document.getElementById('hand');
    handElement.innerHTML = '';
    for (const card of data.Hand) {
        li = document.createElement('li');
        li.textContent = `${card.Suit} ${card.Rank}`;
        li.id = `${card.Id}`;
        handElement.appendChild(li);
        li.addEventListener('click', playCard);
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
}

function playCard(event) {
    const player = parseInt(new URLSearchParams(window.location.search).get('player'));
    
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
    .then(response => response.json())
    .then(data => {
        render(data);
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
            "player": parseInt(new URLSearchParams(window.location.search).get('player'))
        })
    })
    .then(response => response.json())
    .then(data => {
        render(data);
    })
    .catch(error => console.error('Error:', error));
}