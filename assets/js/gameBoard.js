const currentTurn = window.gameState.currentTurn;

document.querySelectorAll('.colonne').forEach(col => {
    col.addEventListener('mouseenter', () => {
        col.querySelectorAll('.cellule').forEach(cell => {
            if (!cell.classList.contains('black') && !cell.classList.contains('orange')) {
                let nb_cell = col.querySelectorAll('.cellule').length - col.querySelectorAll('.black, .orange').length;
                let nb_cell_to_highlight = nb_cell -1;
                let cell_to_hightlight = col.querySelectorAll('.cellule')[nb_cell_to_highlight];
                if (currentTurn === '1'){
                    cell_to_hightlight.classList.add('hover-black','hoverNextTurn');
                } else {
                    cell_to_hightlight.classList.add('hover-orange','hoverNextTurn');
                };
            };
        });
    });

    col.addEventListener('mouseleave', () => {
        col.querySelectorAll('.cellule').forEach(cell => {
            cell.classList.remove('hoverNextTurn')
        })
    });
});