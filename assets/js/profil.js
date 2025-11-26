function createDoughnutChart() {
    // Get canvas element and 2D drawing context
    const canvas = document.getElementById('doughnutChart');
    const ctx = canvas.getContext('2d');
    const size = 200;
    canvas.width = size;
    canvas.height = size;

    // Extract win/loss data from the DOM
    const userWin = parseInt(document.getElementById('victoiresValue').textContent);
    const userLosses = parseInt(document.getElementById('defaitesValue').textContent);
    const total = userWin + userLosses;

    // Handle case where user has no games played
    if (total === 0) {
        ctx.clearRect(0, 0, size, size);
        return;
    }

    // Calculate the angle for the losses segment (in radians)
    const defaitesAngle = (userLosses / total) * 2 * Math.PI;

    // Clear any previous drawing
    ctx.clearRect(0, 0, size, size);

    // Draw losses segment (dark color)
    ctx.beginPath();
    ctx.arc(size/2, size/2, 80, -Math.PI/2, -Math.PI/2 + defaitesAngle);
    ctx.strokeStyle = '#1C1A22';
    ctx.lineWidth = 40;
    ctx.stroke();

    // Draw wins segment (orange color)
    ctx.beginPath();
    ctx.arc(size/2, size/2, 80, -Math.PI/2 + defaitesAngle, -Math.PI/2 + 2 * Math.PI);
    ctx.strokeStyle = '#F28C28';
    ctx.lineWidth = 40;
    ctx.stroke();
}

function createBarChart() {
    // Get canvas element and set dimensions
    const canvas = document.getElementById('barChart');
    const ctx = canvas.getContext('2d');
    const width = 200;
    const height = 150;
    canvas.width = width;
    canvas.height = height;

    // Extract win/loss data from the DOM
    const userWin = document.getElementById('victoiresValue').textContent;
    const userLosses = document.getElementById('defaitesValue').textContent;

    // Calculate scaling based on the maximum value
    const maxValue = Math.max(userWin, userLosses);
    const barWidth = 60;
    const spacing = 40;

    // Clear any previous drawing
    ctx.clearRect(0, 0, width, height);

    // Draw wins bar (orange, left side)
    const victoiresHeight = (userWin / maxValue) * 100;
    ctx.fillStyle = '#F28C28';
    ctx.fillRect(spacing, height - victoiresHeight - 20, barWidth, victoiresHeight);

    // Draw losses bar (dark, right side)
    const defaitesHeight = (userLosses / maxValue) * 100;
    ctx.fillStyle = '#1C1A22';
    ctx.fillRect(spacing + barWidth + 20, height - defaitesHeight - 20, barWidth, defaitesHeight);

    // Draw labels below each bar
    ctx.fillStyle = '#ffffffff';
    ctx.font = '12px Arial';
    ctx.textAlign = 'center';
    ctx.fillText('Victoires', spacing + barWidth/2, height - 5);
    ctx.fillText('Défaites', spacing + barWidth + 20 + barWidth/2, height - 5);
}

document.addEventListener('DOMContentLoaded', () => {
    // Generate the statistical charts
    createDoughnutChart();
    createBarChart();

    // Format the avatar to show first 2 letters of username
    const avatar = document.getElementById('avatar-text');
    avatar.textContent = avatar.textContent.substring(0, 2).toUpperCase();
});