<script lang="ts">
	import LineChart from "$lib/components/LineChart.svelte"

	export let data: {
		analysis: string
		historicalVolumeData: number[],
		historicalSpreadData: number[],
		predictedVolumeData: number[],
		predictedSpreadData: number[],
		xAxis: string[],
		error: string
	}

	$: styledAnalysis = data.analysis
		.replace(/<h2>/g, '<h2 class="text-3xl text-center mb-4">')
		.replace(/<p>/g, '<p class="text-xl mb-2">')
		.replace(/<li>/g, '<li class="text-xl mb-2">')
  </script>

<div class="flex flex-col items-center p-10">
	<h1 class="text-8xl mb-10">STABLETIDE</h1>

	<div class="w-2/3 my-10">
		<p class="">{@html styledAnalysis}</p>
	</div>
  
	<h2 class="text-3xl my-5">Bid-Ask Spread Percentage</h2>
	<div class="mb-10 border-2 border-dotted rounded-lg px-4 py-2 w-full">
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
		  label: 'Bid-ask Spread / Bid-price ',
		}}
	  ></LineChart>
	</div>
  
  <h2 class="text-3xl my-5">Trading Volume</h2>
  <div class="mb-10 border-2 border-dotted rounded-lg px-4 py-2 w-full">
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
  
	{#if data.error}
	  <h2>Error:</h2>
	  <p>{data.error}</p>
	{/if}
</div>

