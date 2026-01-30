document.addEventListener('DOMContentLoaded', function() {
   renderPlayer();
});

function renderPlayer(){
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
        const handElement = document.getElementById('hand');
        handElement.innerHTML = '';
        for (const card of data.Hand) {
            li = document.createElement('li');
            li.textContent = `${card.Suit} ${card.Rank}`;
            li.id = `${card.Id}`;
            handElement.appendChild(li);
            li.addEventListener('click', playCard);
        }
      
      
    })
    .catch(error => console.error('Error:', error));
}

function playCard(event) {
    const player = parseInt(new URLSearchParams(window.location.search).get('player'));
    
    console.log("Player", player, "plays", event.target.id);

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
        renderPlayer();
    })
    .catch(error => console.error('Error:', error));
}