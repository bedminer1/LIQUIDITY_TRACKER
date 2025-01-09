import { readFile } from 'fs/promises'
// @ts-ignore
import path from 'path';

export async function load() {
  const filePath = path.resolve('./src/lib/recommendations.json');

  try {
    const fileData = await readFile(filePath, 'utf-8');
    const recommendations = JSON.parse(fileData)
    let historicalData: LiquidityRecord[] = recommendations.historical_data
    let predictions: LiquidityRecord[] = recommendations.predictions

    let { historicalSpreadData, historicalVolumeData, predictedSpreadData, predictedVolumeData, xAxis } = prepareChartData(historicalData, predictions)

    return { 
        analysis: recommendations.analysis, 
        report: recommendations.report,
        historicalSpreadData: historicalSpreadData,
        historicalVolumeData: historicalVolumeData,
        predictedSpreadData: predictedSpreadData,
        predictedVolumeData: predictedVolumeData,
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
    const historicalSpreadData: (number | null)[] = []
    const predictedSpreadData: (number | null)[] = []
    const historicalVolumeData: (number | null)[] = []
    const predictedVolumeData: (number | null)[] = []
    const xAxis: string[] = []
    for (let record of historicalData) {
        historicalSpreadData.push((record.bid_ask_spread/record.bid_price)*100)
        historicalVolumeData.push(record.volume)
        predictedSpreadData.push(null)
        predictedVolumeData.push(null)
        xAxis.push(record.timestamp.substring(0,10))
    }
    predictedSpreadData[historicalSpreadData.length-1] = historicalSpreadData.at(-1)!
    predictedVolumeData[historicalVolumeData.length-1] = historicalVolumeData.at(-1)!

    for (let record of predictions) {
        predictedSpreadData.push((record.bid_ask_spread/record.bid_price)*100)
        predictedVolumeData.push(record.volume)
        xAxis.push(record.timestamp.substring(0,10))
    }

    return { historicalSpreadData, historicalVolumeData, predictedSpreadData, predictedVolumeData, xAxis };
}
