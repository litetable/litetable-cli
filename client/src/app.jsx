import { Button } from "@/components/ui/button";
import AppLayout from "./layout.jsx";

function App() {
  return (
    <AppLayout>
      <div>
        <h1 className={"text-3xl"}>Welcome to LiteTable!</h1>
        <Button className={"cursor-pointer"}>Click me</Button>
      </div>
    </AppLayout>
  );
}

export default App;
