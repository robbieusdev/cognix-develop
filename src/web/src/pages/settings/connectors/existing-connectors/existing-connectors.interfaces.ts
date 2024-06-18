export namespace Interfaces {
  export interface TableItem {
    id: string;
    source: string;
    // status: LABEL_STATUS | string;
    last_successful_index_time: string | null;
    total_docs_indexed: number;
  }

  export interface UseFilteredHandler {
    columns: ColumnItem[];
    tableData: TableItem[];
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