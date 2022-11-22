import { ThemeProvider } from '@emotion/react'
import { BarChart, ContainerV2, Text, theme } from '@scaleway/ui'
import { useState } from 'react'
import './App.css'
import Drawing from './Drawing'

const HEADER = `
Draw some digits in the canvas below and see what the model think they are!
`;

function App() {
  const [chartData, setChartData] = useState<any[]>([...Array(10).keys()].map(i => ({ id: i, value: 1.0 })));

  const inferDigit = (data: Number[]) => {
    fetch("https://rustexamplesdetfo3u5-rust-mnist-example.functions.fnc.fr-par.scw.cloud", {
      method: "POST",
      body: JSON.stringify({
        data: data
      }),
    }).then(resp => resp.json()).then((resp: { output: number[] }) => {
      let pos = resp.output.map(val => Math.max(0, val))
      const sumPos = pos.reduce(
        (accumulator, currentValue) => accumulator + currentValue,
        0
      );

      setChartData(pos.map((val, index) => ({ id: index, value: (100 * (val / sumPos)).toFixed(1) })));
    })
  };

  return (
    <ThemeProvider theme={theme}>
      <ContainerV2 title='Serverless MNIST with Rust' header={HEADER} className="App">
        <Drawing inferDigit={inferDigit}></Drawing>
        <BarChart data={chartData}></BarChart>
      </ContainerV2>
    </ThemeProvider>
  )
}

export default App
