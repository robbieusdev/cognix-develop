import { LABEL_STATUS, LabelStatus } from "@/components/ui/label-status";
import { TableRow, TableCell, TableBody, Table } from "@/components/ui/table";

import PauseIcon from "@/assets/svgs/pause-icon.svg?react";
import EditIcon from "@/assets/svgs/edit-icon.svg?react";
import TrashIcon from "@/assets/svgs/trash-icon.svg?react";

import { memo } from "react";
import { RenderTableHeader } from "./render-table-header";

interface TableItem {
  id: string;
  connector: string;
  status: LABEL_STATUS | string;
  last_indexed: string;
  docs_indexed: string;
}

interface ColumnItem {
  label: string;
  accessor: string;
}

interface Props {
  tableData: any[];
  columns: ColumnItem[];
  sortField: string;
  handleSortingChange: (value: string) => void;
  onPause?: (value: string) => void;
  onEdit: (value: string) => void;
  onDelete: (value: string) => void;
  withBtn?: boolean;
}

const RenderTable = memo(
  ({
    tableData,
    columns,
    sortField,
    handleSortingChange,
    onPause,
    onEdit,
    onDelete,
    withBtn,
  }: Props) => {
    return (
      <Table className="bg-main">
        <RenderTableHeader
          columns={columns}
          handleSortingChange={handleSortingChange}
          sortField={sortField}
        />
        <TableBody>
          {tableData.map((data: TableItem) => {
            return (
              <TableRow key={data.id}>
                {columns.map(({ accessor }: ColumnItem) => {
                  const tData = (data as Record<string, any>)[accessor] ?? "——";
                  return (
                    <TableCell key={accessor}>
                      <LabelStatus label={tData} />
                    </TableCell>
                  );
                })}
                {withBtn && (
                  <TableCell>
                    <div className="flex flex-row items-center justify-center gap-3">
                      {onPause && (
                        <PauseIcon
                          onClick={() => onPause(data.id)}
                          className="cursor-pointer"
                        />
                      )}
                      <EditIcon
                        onClick={() => onEdit(data.id)}
                        className="cursor-pointer"
                      />
                      <TrashIcon
                        onClick={() => onDelete(data.id)}
                        className="cursor-pointer"
                      />
                    </div>
                  </TableCell>
                )}
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    );
  }
);

export { RenderTable };
