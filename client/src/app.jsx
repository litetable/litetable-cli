import AppLayout from "./layout";
import QueryInput from "@/components/query-input";

function App() {
  return (
    <AppLayout>
      <div>
        <h1 className={"text-3xl"}>Welcome to LiteTable!</h1>
        <QueryInput />
      </div>
    </AppLayout>
  );
}

export default App;
