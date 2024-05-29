import { MetaFunction, json, ActionFunction } from "@remix-run/node";
import { getUser } from "../grpc/user.server";
import { useLoaderData, useFetcher } from "@remix-run/react";

export const meta: MetaFunction = () => {
  return [
    { title: "New Remix App" },
    { name: "description", content: "Welcome to Remix!" },
  ];
};

export const loader = async () => {
  const user = await getUser(1); // Default user ID 1
  return json({ ...user });
};

export const action: ActionFunction = async ({ request }) => {
  const formData = await request.formData();
  const userId = formData.get("userId")?.toString() || "1";
  const user = await getUser(parseInt(userId, 10));
  return json({ ...user });
};

export default function Index() {
  const initialData = useLoaderData();
  const fetcher = useFetcher();

  // Check if fetcher has data, otherwise use initial data
  const user = fetcher.data || initialData;

  return (
    <div style={{ fontFamily: "system-ui, sans-serif", lineHeight: "1.8" }}>
      <h1>User</h1>
      <fetcher.Form method="post">
        <input type="number" name="userId" defaultValue={1} />
        <button type="submit">Fetch User</button>
      </fetcher.Form>
      <pre>{JSON.stringify(user, null, 2)}</pre>
    </div>
  );
}
