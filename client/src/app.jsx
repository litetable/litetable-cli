import { HashRouter, Routes, Route } from "react-router-dom";
import AppLayout from "./layout";
import QueryBuilder from "@/components/query-builder";
import JsonUpload from "@/components/json-upload";

function Home() {
  return (
    <div>
      <h1 className="text-3xl font-bold text-center mb-6">
        Welcome to LiteTable!
      </h1>
    </div>
  );
}

function App() {
  return (
    <HashRouter>
      <AppLayout>
        <div className="container mx-auto max-w-[1200px] py-8">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/query" element={<QueryBuilder />} />
            <Route path="/upload" element={<JsonUpload />} />
          </Routes>
        </div>
      </AppLayout>
    </HashRouter>
  );
}

export default App;
