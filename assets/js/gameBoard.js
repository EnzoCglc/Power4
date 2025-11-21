let currentTurn = window.gameState.currentTurn;
let player1 = window.gameState.player1;
let Finish = false;

// Fonction pour mettre à jour l'affichage du joueur actif
function updatePlayerIndicator(turn) {
    const turnLayer1 = document.querySelector('.turn-hud-layer1');
    const turnLayer2 = document.querySelector('.turn-hud-layer2');

    if (turn === 1) {
        turnLayer1.style.display = 'flex';
        turnLayer2.style.display = 'none';
    } else {
        turnLayer1.style.display = 'none';
        turnLayer2.style.display = 'flex';
    }
}

// Initialiser l'affichage au chargement
updatePlayerIndicator(currentTurn);

document.querySelectorAll('.colonne').forEach(col => {
    col.addEventListener('mouseenter', () => {
        col.querySelectorAll('.cellule').forEach(cell => {
            if (!cell.classList.contains('black') && !cell.classList.contains('orange')) {
                let nb_cell = col.querySelectorAll('.cellule').length - col.querySelectorAll('.black, .orange').length;
                let nb_cell_to_highlight = nb_cell -1;
                let cell_to_hightlight = col.querySelectorAll('.cellule')[nb_cell_to_highlight];
                if (Finish === true){
                    cell.classList.remove('hoverNextTurn', 'hover-black', 'hover-orange');
                } else {
                    if (currentTurn === 1){
                        cell_to_hightlight.classList.add('hover-orange','hoverNextTurn');
                    } else {
                        cell_to_hightlight.classList.add('hover-black','hoverNextTurn');
                    };
                }
            };
        });
    });

    col.addEventListener('mouseleave', () => {
        col.querySelectorAll('.cellule').forEach(cell => {
            cell.classList.remove('hoverNextTurn', 'hover-black', 'hover-orange');
        })
    });
});

function playColumn(colIndex) {
   
    fetch('/game', {
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
        Cell.classList.add('orange', 'animate-drop');
    } else {
        Cell.classList.add('black', 'animate-drop');
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
    window.gameState.currentTurn = game.CurrentTurn;
    currentTurn = game.CurrentTurn;
    Finish = game.GameOver;

    // Mettre à jour l'indicateur du joueur actif
    updatePlayerIndicator(currentTurn);

    document.querySelectorAll('.cellule').forEach(cell => {
        cell.classList.remove('hoverNextTurn', 'hover-black', 'hover-orange');
    });
    if (Finish === true) {
        const winMsg = document.getElementById('win-msg');
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
        
        if (game.Winner === 1) {
            winMsg.textContent = `${player1} Win a game`
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
