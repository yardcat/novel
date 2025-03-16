import React, { useRef, useEffect, useState } from 'react';

const GridCanvas = () => {
  const canvasRef = useRef(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const gridSize = 1;
  const rows = 50;
  const cols = 50;

  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');

    const drawGrid = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      for (let i = 0; i < rows; i++) {
        for (let j = 0; j < cols; j++) {
          ctx.strokeStyle = '#000';
          ctx.strokeRect(j * gridSize, i * gridSize, gridSize, gridSize);
        }
      }
    };

    drawGrid();
  }, []);

  const fillCell = (x, y) => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    ctx.fillStyle = '#000';
    ctx.fillRect(x * gridSize, y * gridSize, gridSize, gridSize);
  };

  const getCellCoordinates = (x, y) => {
    const cellX = Math.floor(x / gridSize);
    const cellY = Math.floor(y / gridSize);
    return { cellX, cellY };
  };

  const handleMouseDown = (e) => {
    setIsDrawing(true);
    const { offsetX, offsetY } = e.nativeEvent;
    const { cellX, cellY } = getCellCoordinates(offsetX, offsetY);
    fillCell(cellX, cellY);
  };

  const handleMouseMove = (e) => {
    if (isDrawing) {
      const { offsetX, offsetY } = e.nativeEvent;
      const { cellX, cellY } = getCellCoordinates(offsetX, offsetY);
      fillCell(cellX, cellY);
    }
  };

  const handleMouseUp = () => {
    setIsDrawing(false);
  };

  const handleMouseLeave = () => {
    setIsDrawing(false);
  };

  return (
    <canvas
      ref={canvasRef}
      width={400}
      height={400}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseLeave}
      style={{ border: '1px solid #000', cursor: 'pointer' }}
    />
  );
};

export { GridCanvas };
