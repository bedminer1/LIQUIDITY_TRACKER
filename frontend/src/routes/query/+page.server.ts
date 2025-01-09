import { saveResponseToFile } from '$lib/utils/saveResponse'
import { redirect } from '@sveltejs/kit';

export const actions = {
  default: async ({ request }) => {
    const formData = await request.formData();
    const start = formData.get("start");
    const end = formData.get("end");
    const asset = formData.get("asset");
    const time_intervals = formData.get("time_intervals");
    const time_interval_length = formData.get("time_interval_length");

    // Validate inputs
    if (!start || !end || !asset || !time_intervals || !time_interval_length) {
      return { error: "All fields are required." };
    }

    const apiUrl = `http://localhost:4000/recommendations?start=${start}&end=${end}&asset=${asset}&time_intervals=${time_intervals}&time_interval_length=${time_interval_length}`;

    try {
      // Fetch data from the backend
      const response = await fetch(apiUrl);

      if (!response.ok) {
        throw new Error(`API request failed with status ${response.status}`);
      }

      const data = await response.json();
      await saveResponseToFile(data, "recommendations.json");

      // Return the API response to the frontend
      console.log("redirecting..")
      redirect(303, "/")
    } catch (error) {
      return { error: "Something went wrong" };
    }
  },
};
