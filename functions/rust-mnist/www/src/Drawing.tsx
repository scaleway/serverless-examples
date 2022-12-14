import { MouseEvent, Ref, useEffect, useRef, useState } from 'react';
import { theme, normalize, Button } from '@scaleway/ui'

const MNIST_SIZE = 28;

export default function Drawing(props: { inferDigit: (data: Array<Number>) => void }) {

    const [drawing, setDrawing] = useState(false);
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const scaledCanvasRef = useRef<HTMLCanvasElement>(null);
    const ctxRef = useRef<null | CanvasRenderingContext2D>(null);
    const scaledCtxRef = useRef<null | CanvasRenderingContext2D>(null);

    const startDraw = (e: MouseEvent) => {
        if (!ctxRef.current) {
            return;
        }
        const { offsetX, offsetY } = e.nativeEvent;
        ctxRef.current.beginPath();
        ctxRef.current.moveTo(offsetX, offsetY);
        setDrawing(true);

    };

    const stopDraw = () => {
        if (!ctxRef.current) {
            return;
        }
        ctxRef.current.closePath();
        setDrawing(false);

        if (!scaledCtxRef.current || !canvasRef.current) {
            return;
        }

        let ctxScaled = scaledCtxRef.current;
        ctxScaled.save();
        ctxScaled.clearRect(0, 0, ctxScaled.canvas.height, ctxScaled.canvas.width);
        ctxScaled.scale(MNIST_SIZE / canvasRef.current.width, MNIST_SIZE / canvasRef.current.height)
        ctxScaled.drawImage(canvasRef.current, 0, 0)
        const { data } = ctxScaled.getImageData(0, 0, MNIST_SIZE, MNIST_SIZE)
        ctxScaled.restore();

        let pixelData = new Array(MNIST_SIZE * MNIST_SIZE).fill(0.0);
        for (let i = 0; i < Math.floor(data.length / 4); i++) {
            let pixel = { r: data[4 * i], g: data[4 * i + 1], b: data[4 * i + 2] };
            if (pixel.r == 57.0 && pixel.g == 1.0 && pixel.b == 113.0) {
                pixelData[i] = 1.0;
            }
        }

        props.inferDigit(pixelData);
    };

    const draw = (e: MouseEvent) => {
        if (!ctxRef.current || !drawing) {
            return;
        }
        const { offsetX, offsetY } = e.nativeEvent;
        ctxRef.current.lineTo(offsetX, offsetY);
        ctxRef.current.stroke();
    };

    const clear = () => {
        if (!ctxRef.current || !canvasRef.current) {
            return;
        }
        ctxRef.current.clearRect(0, 0, canvasRef.current.width, canvasRef.current.height);
        ctxRef.current.fillStyle = theme.colors.primary.background;
        ctxRef.current.fillRect(0, 0, canvasRef.current.width, canvasRef.current.height);
        setDrawing(false);
    };

    useEffect(() => {
        const canvas = canvasRef.current;
        if (!canvas) {
            return;
        }

        canvas.style.width = '50%';
        canvas.style.height = `${canvas.offsetWidth}px`;
        canvas.width = 2 * canvas.offsetWidth;
        canvas.height = 2 * canvas.offsetHeight;

        // Setting the context to enable us draw
        const ctx = canvas.getContext('2d');
        if (!ctx) {
            return;
        }

        ctx.scale(2, 2);
        ctx.lineCap = 'round';
        ctx.strokeStyle = theme.colors.primary.text;
        ctx.lineWidth = 20;

        ctx.fillStyle = theme.colors.primary.background;
        ctx.fillRect(0, 0, canvas.width, canvas.height);

        ctxRef.current = ctx;

        const scaledCtx = scaledCanvasRef.current?.getContext('2d');
        if (!scaledCtx) {
            return;
        }

        scaledCtxRef.current = scaledCtx;


    }, []);



    return (
        <div style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
            <canvas id="digit-canvas"
                onMouseDown={startDraw}
                onMouseUp={stopDraw}
                onMouseMove={draw}
                ref={canvasRef}
            />
            <canvas ref={scaledCanvasRef} style={{ height: `${MNIST_SIZE}px`, width: `${MNIST_SIZE}px` }}></canvas>
            <Button onClick={clear}>Clear</Button>
        </div >
    );

}