<script lang="ts">
    import { enhance } from "$app/forms";
    let start = "2016-01-01";
	let end = "2018-01-01";
	let asset = "";
	let time_intervals = 50;
	let time_interval_length = 86400;
</script>

<div class="flex flex-col items-center p-10">
    <h2 class="text-4xl mb-10">QUERY</h2>
    <form method="post" use:enhance class="flex flex-col gap-5 w-1/2 text-xl items-center border-2 border-dotted rounded-lg p-5 mb-10">
        <label class="w-full flex justify-between">
          <p class="flex items-center text-xl">
              Start Date:
          </p>
          <input type="date" class="input w-2/3" name="start" bind:value={start} required />
        </label>
      
        <label class="w-full flex justify-between">
          <p class="flex items-center text-xl">
              End Date:
          </p>
          <input type="date" class="input w-2/3" name="end" bind:value={end} required />
        </label>
      
        <label class="w-full flex justify-between">
          <p class="flex items-center text-xl">
              Asset:
          </p>
          <input type="text" class="input w-2/3" name="asset" bind:value={asset} placeholder="e.g., ETF_EMB" required />
        </label> 
      
        <label class="w-full flex justify-between">
          <p class="flex items-center text-xl">
              Time Intervals:
          </p>
          <input type="number" class="input w-2/3" name="time_intervals" bind:value={time_intervals} min="1" required />
        </label>
      
        <label class="w-full flex justify-between mb-4">
          <p class="flex items-center text-xl">
              Time Interval:
          </p>
          <select name="time_interval_length" class="select w-2/3" bind:value={time_interval_length}>
              <option value="3600">Hour</option>
              <option value="86400">Day</option>
              <option value="604800">Week</option>
              <option value="18144000">Month</option>
          </select>
        </label>
      
        <button class="btn  text-lg variant-ghost-primary" type="submit">Submit</button>
      </form>
</div>