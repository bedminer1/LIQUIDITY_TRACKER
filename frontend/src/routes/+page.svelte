<script lang="ts">
	import LineChart from "$lib/components/LineChart.svelte"
	import Card from "$lib/components/Card.svelte";
	import { onMount } from "svelte"

	export let data: {
		analysis: string,
		report: LiquidityReport,
		historicalVolumeData: number[],
		historicalSpreadData: number[],
		predictedVolumeData: number[],
		predictedSpreadData: number[],
		xAxis: string[],
		error: string,
		currentDay: string,
	}

	$: styledAnalysis = data.analysis
		.replace(/<h2>/g, '<h2 class="text-3xl text-center mb-4">')
		.replace(/<p>/g, '<p class="text-xl mb-2">')
		.replace(/<li>/g, '<li class="text-xl mb-2">')
	let historicalIlliquidityRate = ((data.report.current_moderate_risk_count+data.report.current_high_risk_count)/data.report.historical_records).toFixed(3)
	let predictedIlliquidityRate = ((data.report.predicted_moderate_risk_count+data.report.predicted_high_risk_count)/data.report.prediction_records).toFixed(3)

	function formatDate(dateString: string): string {
		const months = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"];
		const [year, month, day] = dateString.split("-");
		const monthShort = months[parseInt(month, 10) - 1];
		return `${parseInt(day, 10)} ${monthShort} ${year}`;
	}

	let time = ""
	function updateClock() {
		const now = new Date();
		const hours = now.getHours().toString().padStart(2, "0");
		const minutes = now.getMinutes().toString().padStart(2, "0");
		const seconds = now.getSeconds().toString().padStart(2, "0");
		time = `${hours}:${minutes}:${seconds}`;
	}
	$: displayDate = time + " " + formatDate(data.currentDay)
	let interval: ReturnType<typeof setInterval>;
	onMount(() => {
		updateClock(); // Initialize the clock immediately
		interval = setInterval(updateClock, 1000); // Update every second
		return () => clearInterval(interval); // Cleanup on component destroy
	});
  </script>

<div class="flex flex-col items-center p-10">
	<h1 class="text-8xl mb-3"><a href="https://github.com/bedminer1/LIQUIDITY_TRACKER/">STABLETIDE</a></h1>
	<p class="mb-8">{displayDate}</p>
	<nav class="underline text-left w-full mb-4">
		<a href="/query">New Query?</a>
	</nav>
	<!-- GRAPH AND ANALYSIS -->
	<div class="flex w-full gap-4 items-start h-[100vh]">
		<!-- ANALYSIS -->
		<div class="w-5/12 justify-start flex flex-col gap-4 h-full">
			<div class="flex flex-col gap-4">
				<Card
					{...{
						title: "Asset",
						body: data.report.asset_type,
						subtitle: "Represents 1/10 of the ETF by JP Morgan",
						icon: "&#9814;"
					}}
				/>
				<Card
					{...{
						title: "Predicted Illiquidity Rate",
						body: predictedIlliquidityRate + "%",
						subtitle: (predictedIlliquidityRate < historicalIlliquidityRate ? "down" : "up") + " compared to " + historicalIlliquidityRate + "% historically",
						icon: "&#10150;"
					}}
				/>
				<Card
					{...{
						title: "Period Analyzed",
						body: formatDate(data.xAxis[0]) + " to " + formatDate(data.currentDay),
						subtitle: "Predicting from " + formatDate(data.currentDay) + " to " + formatDate(data.xAxis.at(-1)!),
						icon: "&#9790;"
					}}
				/>
			</div>
			<p class="card p-4">{@html styledAnalysis}</p>
		</div>
		<!-- GRAPH -->
		<div class="flex flex-col gap-4 w-7/12">
			<div class="card p-6">
				<h2 class="text-3xl mb-3 text-center">Bid-Ask Spread Percentage</h2>
				<div class="mb-3 border-2 border-dotted h-[30vh] rounded-lg px-4 py-2 w-full flex items-center">
					<LineChart 
					{...{
						stats: [
						{
							label: 'Historical Bid-Ask Spread Percentage',
							data: data.historicalSpreadData,
							xAxis: data.xAxis,
							borderColor: 'rgba(75, 192, 192, 1)',
							backgroundColor: 'rgba(75, 192, 192, 0.2)',
						},
						{
							label: 'Predicted Bid-Ask Spread Percentage',
							data: data.predictedSpreadData,
							xAxis: data.xAxis,
							borderColor: 'rgba(255, 99, 132, 1)',
							backgroundColor: 'rgba(255, 99, 132, 0.2)',
						},
						],
						label: 'Bid-ask Spread / Bid-price',
					}}
					></LineChart>
				</div>
			</div>
			<div class="card p-6">
				<h2 class="text-3xl mb-5 text-center">Trading Volume</h2>
				<div class="mb-3 border-2 border-dotted h-[30vh] rounded-lg px-4 py-2 w-full flex items-center">
					<LineChart
					{...{
						stats: [
						{
							label: 'Historical Trading Volume Percentage',
							data: data.historicalVolumeData,
							xAxis: data.xAxis,
							borderColor: 'rgba(75, 192, 192, 1)',
							backgroundColor: 'rgba(75, 192, 192, 0.2)',
						},
						{
							label: 'Predicted Trading Volume Percentage',
							data: data.predictedVolumeData,
							xAxis: data.xAxis,
							borderColor: 'rgba(255, 99, 132, 1)',
							backgroundColor: 'rgba(255, 99, 132, 0.2)',
						},
						],
						label: 'Trading Volume',
					}}
					></LineChart>
				</div>
			</div>
		</div>
	</div>
	</div>
	<div>
	
	
	{#if data.error}
		<h2>Error:</h2>
		<p>{data.error}</p>
	{/if}
</div>

