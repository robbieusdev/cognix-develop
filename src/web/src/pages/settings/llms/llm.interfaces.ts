import { Persona } from "@/models/settings";

export namespace Interfaces {

  export interface UseFilteredHandler {
    columns: ColumnItem[];
    tableData: Persona[];
    sortField: string;
    handleSortingChange: (accessor: string) => void;
  }

  export interface ColumnItem {
    label: string;
    accessor: string;
  }

  export interface Tabs {
    key: string;
    label: string;
  }
}