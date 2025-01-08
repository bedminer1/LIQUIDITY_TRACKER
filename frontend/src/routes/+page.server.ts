import { saveResponseToFile } from '$lib/utils/saveResponse'
import { readFile } from 'fs/promises';
import path from 'path';

export async function load() {
  const filePath = path.resolve('./src/lib/recommendations.json');

  try {
    const fileData = await readFile(filePath, 'utf-8');
    const recommendations = JSON.parse(fileData)
    let historicalData: LiquidityRecord[] = recommendations.historical_data
    let predictions: LiquidityRecord[] = recommendations.predictions

    let { spreadData, volumeData, xAxis } = prepareChartData(historicalData, predictions)

    return { 
        analysis: recommendations.analysis, 
        spreadData: spreadData,
        volumeData: volumeData,
        xAxis: xAxis,
       
    }
  } catch (error) {
    console.error('Error reading recommendations file:', error);
    return { 
        analysis: null,
        spreadData: null,
        volumeData: null,
        xAxis: null,
    }
  }
}

function prepareChartData(historicalData: LiquidityRecord[], predictions: LiquidityRecord[]) {
    const combinedData = [...historicalData, ...predictions];

    const spreadData: number[] = []
    const volumeData: number[] = []
    const xAxis: string[] = []
    for (let record of combinedData) {
        spreadData.push((record.bid_ask_spread/record.bid_price)*100)
        volumeData.push(record.volume)
        xAxis.push(record.timestamp.substring(0,10))
    }

    return { spreadData, volumeData, xAxis };
}
