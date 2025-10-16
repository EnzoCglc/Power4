let userData = {
    username: 'enzo',
    elo: 1020,
    victoires: 50,
    defaites: 40,
    historique: [
        { date: '2025-10-14', adversaire: 'Player1', resultat: 'victoire', eloChange: +15 },
        { date: '2025-10-13', adversaire: 'Player2', resultat: 'defaite', eloChange: -12 },
        { date: '2025-10-13', adversaire: 'Player3', resultat: 'victoire', eloChange: +18 },
        { date: '2025-10-12', adversaire: 'Player4', resultat: 'victoire', eloChange: +14 },
        { date: '2025-10-11', adversaire: 'Player5', resultat: 'defaite', eloChange: -10 },
        { date: '2025-10-10', adversaire: 'Player6', resultat: 'victoire', eloChange: +16 },
        { date: '2025-10-09', adversaire: 'Player7', resultat: 'defaite', eloChange: -13 },
        { date: '2025-10-08', adversaire: 'Player8', resultat: 'victoire', eloChange: +17 }
    ]
};

function updateDisplay() {
        document.getElement
        document.getElementById('displayUsername').textContent = userData.username;
        document.getElementById('usernameDisplay').textContent = userData.username;
        document.getElementById('eloValue').textContent = userData.elo;
        document.getElementById('eloMain').textContent = userData.elo;
        document.getElementById('victoiresValue').textContent = userData.victoires;
        document.getElementById('defaitesValue').textContent = userData.defaites;
  
    const total = userData.victoires + userData.defaites;
    const winrate = Math.round((userData.victoires / total) * 100);
    document.getElementById('winrateText').textContent = `${winrate}% de victoires`;
}

