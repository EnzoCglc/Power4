let currentTurn = window.gameState.currentTurn;
let Finish = false;

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
                    if (window.gameState.currentTurn === 1){
                        console.log('test')
                        cell_to_hightlight.classList.add('hover-black','hoverNextTurn');
                    } else {
                        console.log('test')
                        cell_to_hightlight.classList.add('hover-orange','hoverNextTurn');
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
    console.log('Envoi du coup pour la colonne:', colIndex);
    
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
        console.log('✅ Réponse reçue du backend:', data);
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

function updateGrid(game) {
    console.log('Update du Grid');

    game.Columns.forEach((col, colIndex) => {
        const Col = document.querySelectorAll('.colonne')[colIndex];

        col.forEach((cell, rowIndex) => {
            const Cell = Col.querySelectorAll('.cellule')[rowIndex];

           // console.log("Value de la cellule : ", Cell)

            if (cell === 1) {
                Cell.classList.add('black');
                console.log(`Pion du joueur en [${colIndex}][${rowIndex}]`);

            } else if (cell === 2) {
                Cell.classList.add('orange');
                console.log(`Pion du joueur en [${colIndex}][${rowIndex}]`);
            }
        });
    });
    window.gameState.currentTurn = game.CurrenctTurn;
    Finish = game.GameOver;
    console.log("Valeur du json ", game.CurrenctTurn);
    console.log("Type du currentTurn :", typeof window.gameState.currentTurn);
    console.log("Valeur du tour ", window.gameState.currentTurn);

    document.querySelectorAll('.cellule').forEach(cell => {
        cell.classList.remove('hoverNextTurn', 'hover-black', 'hover-orange');
    });
    if (Finish === true) {
        document.querySelector('.win-banner-overlay').style.display = 'flex';
    }
}

document.querySelector('.win-banner-overlay').style.display = 'none';

    
const bg = document.getElementById('bg');
    bg.play().catch(()=>{ /* blocked until user gesture */ });

// on first user gesture unmute and resume playback
function unlockAudio(e){
  bg.muted = false;
  bg.play().catch(err => console.warn('play blocked:', err));
  window.removeEventListener('pointerdown', unlockAudio);
}

window.addEventListener('pointerdown', unlockAudio, { once: true });
