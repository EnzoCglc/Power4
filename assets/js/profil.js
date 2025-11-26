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
    const canvas = document.getElementById('barChart');
    const ctx = canvas.getContext('2d');
    const config = setupBarChartCanvas(canvas);

    const data = getBarChartData();
    ctx.clearRect(0, 0, config.width, config.height);

    drawBars(ctx, data, config);
    drawBarLabels(ctx, config);
}

// setupBarChartCanvas initializes canvas dimensions.
function setupBarChartCanvas(canvas) {
    const width = 200, height = 150;
    canvas.width = width;
    canvas.height = height;
    return { width, height, barWidth: 60, spacing: 40 };
}

// getBarChartData extracts win/loss data from DOM.
function getBarChartData() {
    const userWin = parseInt(document.getElementById('victoiresValue').textContent);
    const userLosses = parseInt(document.getElementById('defaitesValue').textContent);
    const maxValue = Math.max(userWin, userLosses);
    return { userWin, userLosses, maxValue };
}

// drawBars draws the win/loss bars on the canvas.
function drawBars(ctx, data, config) {
    const victoiresHeight = (data.userWin / data.maxValue) * 100;
    ctx.fillStyle = '#F28C28';
    ctx.fillRect(config.spacing, config.height - victoiresHeight - 20, config.barWidth, victoiresHeight);

    const defaitesHeight = (data.userLosses / data.maxValue) * 100;
    ctx.fillStyle = '#1C1A22';
    ctx.fillRect(config.spacing + config.barWidth + 20, config.height - defaitesHeight - 20, config.barWidth, defaitesHeight);
}

// drawBarLabels adds text labels below bars.
function drawBarLabels(ctx, config) {
    ctx.fillStyle = '#ffffffff';
    ctx.font = '12px Arial';
    ctx.textAlign = 'center';
    ctx.fillText('Victoires', config.spacing + config.barWidth/2, config.height - 5);
    ctx.fillText('Défaites', config.spacing + config.barWidth + 20 + config.barWidth/2, config.height - 5);
}

document.addEventListener('DOMContentLoaded', () => {
    // Generate the statistical charts
    createDoughnutChart();
    createBarChart();

    // Format the avatar to show first 2 letters of username
    const avatar = document.getElementById('avatar-text');
    avatar.textContent = avatar.textContent.substring(0, 2).toUpperCase();
});