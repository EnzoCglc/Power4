let currentTurn = window.gameState.currentTurn;
let player1 = window.gameState.player1;
let gameMode = window.gameState.gameMode || "duo";
let botLevel = window.gameState.botLevel || 0;
let Finish = false;

document.querySelectorAll('.colonne').forEach(col => {
    col.addEventListener('mouseenter', () => {
        if (Finish === true) return;

        // Trouver la première cellule disponible (la plus basse)
        let nb_cell = col.querySelectorAll('.cellule').length - col.querySelectorAll('.black, .orange').length;
        let nb_cell_to_highlight = nb_cell - 1;

        if (nb_cell_to_highlight >= 0) {
            let cell_to_highlight = col.querySelectorAll('.cellule')[nb_cell_to_highlight];

            if (currentTurn === 1) {
                cell_to_highlight.classList.add('hover-black', 'hoverNextTurn');
            } else {
                cell_to_highlight.classList.add('hover-orange', 'hoverNextTurn');
            }
        }
    });

    col.addEventListener('mouseleave', () => {
        col.querySelectorAll('.cellule').forEach(cell => {
            cell.classList.remove('hoverNextTurn', 'hover-black', 'hover-orange');
        })
    });
});

function playColumn(colIndex) {
    const endpoint = gameMode === "bot" ? '/game/bot/play' : '/game';

    fetch(endpoint, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            col: colIndex
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            updateGrid(data.data.game);
        } else {
            console.error('Erreur:', error);
        }
    })
    .catch(error => {
        console.error('Erreur:', error);
    });
}

function dropToken(colIndex, rowIndex, player) {
    const Col = document.querySelectorAll('.colonne')[colIndex];
    const Cell = Col.querySelectorAll('.cellule')[rowIndex];

    Cell.classList.remove('black', 'orange', 'animate-drop');

    if (player === 1) {
        Cell.classList.add('black', 'animate-drop');
    } else {
        Cell.classList.add('orange', 'animate-drop');
    }

    Cell.addEventListener('animationEnd', () => {
        Cell.classList.remove('animate-drop');
    })
}

function updateGrid(game) {
    game.Columns.forEach((col, colIndex) => {
        const Col = document.querySelectorAll('.colonne')[colIndex];

        col.forEach((cell, rowIndex) => {
            const Cell = Col.querySelectorAll('.cellule')[rowIndex];

            if (cell === 1 && !Cell.classList.contains('black')) {
                dropToken(colIndex, rowIndex, 1);
            } else if (cell === 2 && !Cell.classList.contains('orange')) {
                dropToken(colIndex, rowIndex, 2);
            }
        });
    });
    currentTurn = game.CurrentTurn;
    window.gameState.currentTurn = game.CurrentTurn;
    Finish = game.GameOver;

    document.querySelectorAll('.cellule').forEach(cell => {
        cell.classList.remove('hoverNextTurn', 'hover-black', 'hover-orange');
    });
    if (Finish === true) {
        const winMsg = document.getElementById('win-msg');

        // Ne sauvegarder le résultat que pour le mode duo
        if (gameMode === "duo") {
            const body = {
                winner: game.Winner,
                player1: player1,
                player2: "player2",
                isDraw: game.isDraw
            };

            fetch('/game/result', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body)
            })
            .then(res => res.json())
            .then(data => console.log("Résultat enregistré:", data))
            .catch(err => console.error("Erreur update ELO:", err));
        }

        // Message adapté selon le mode
        if (game.Winner === 1) {
            winMsg.textContent = `${player1} Win a game`
        } else if (gameMode === "bot") {
            winMsg.textContent = `🤖 Bot Level ${botLevel} Win a game`
        } else {
            winMsg.textContent = "Player 2 Win a game"
        }

        document.querySelector('.win-banner-overlay').style.display = 'flex';

    }
}

document.querySelector('.win-banner-overlay').style.display = 'none';

    
const bg = document.getElementById('bg');

function unlockAudio() {
  if (bg.muted) {
    bg.muted = false;
  }
  bg.play().catch(err => console.warn('play blocked:', err));
}

window.addEventListener('pointerdown', unlockAudio, { once: true });
