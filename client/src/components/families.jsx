import { useState, useEffect } from "react";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
  DialogClose,
} from "@/components/ui/dialog";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";
import { useQueryApi } from "@/hooks/use-query-api";

export default function FamilySelector({ columnFamily, setColumnFamily }) {
  const [newColumnFamily, setNewColumnFamily] = useState("");
  const [families, setFamilies] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const { handleSubmit } = useQueryApi();

  // Fetch families from API on component mount
  useEffect(() => {
    const fetch = async () => {
      await fetchFamilies();
    };
    fetch();
  }, []);

  const fetchFamilies = async () => {
    setIsLoading(true);
    try {
      const response = await fetch("/families");
      if (!response.ok) {
        throw new Error("Failed to fetch families");
      }
      const data = await response.json();
      setFamilies(data);

      // Set default selection if there's no current selection
      if (!columnFamily && data.length > 0) {
        setColumnFamily(data[0]);
      }
    } catch (err) {
      setError(err.message);
      toast.error("Failed to load column families");
    } finally {
      setIsLoading(false);
    }
  };

  const addNewColumnFamily = async () => {
    const newest = newColumnFamily.trim();
    if (!newest || families.includes(newest)) {
      toast.error("Family name must be unique and not empty");
      return;
    }

    try {
      const famQuery = {
        type: "CREATE",
        families: [newest],
      };

      await handleSubmit(famQuery);

      setNewColumnFamily("");
      toast("Column family created");

      // ✅ Refresh from backend
      await fetchFamilies();

      // ✅ Update the selection
      setColumnFamily(newest);
    } catch (e) {
      toast.error("Failed to create column family");
      console.error(e);
    }
  };

  if (error) {
    return <div className="text-red-500">Error loading families: {error}</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-1">
        <Label htmlFor="columnFamily">Column Family</Label>
        <Dialog>
          <DialogTrigger asChild>
            <Button variant="ghost" size="sm" className="h-6 px-2">
              <Plus className="h-3.5 w-3.5 mr-1" />
              Add New
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Add New Column Family</DialogTitle>
            </DialogHeader>
            <div className="py-4">
              <Label htmlFor="newColumnFamily" className="mb-2 block">
                Name
              </Label>
              <Input
                id="newColumnFamily"
                value={newColumnFamily}
                onChange={(e) => setNewColumnFamily(e.target.value)}
                placeholder="Enter column family name"
              />
            </div>
            <DialogFooter>
              <DialogClose asChild>
                <Button variant="outline">Cancel</Button>
              </DialogClose>
              <DialogClose asChild>
                <Button
                  onClick={addNewColumnFamily}
                  disabled={!newColumnFamily}
                >
                  Add
                </Button>
              </DialogClose>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
      <Select
        value={columnFamily}
        onValueChange={setColumnFamily}
        disabled={isLoading}
      >
        <SelectTrigger id="columnFamily">
          {isLoading ? (
            <div className="flex items-center">
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              Loading...
            </div>
          ) : (
            <SelectValue placeholder="Select column family" />
          )}
        </SelectTrigger>
        <SelectContent>
          {families.map((family) => (
            <SelectItem key={family} value={family}>
              {family}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
