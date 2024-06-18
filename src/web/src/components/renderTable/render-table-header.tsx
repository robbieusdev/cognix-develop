import ArrowActive from '@/assets/svgs/filter-arrow-active.svg?react';
import ArrowUnActive from '@/assets/svgs/filter-arrow-unactive.svg?react';
import { TableCell, TableHeader, TableRow } from '@/components/ui/table';
import { memo } from 'react';

interface ColumnItem {
  label: string;
  accessor: string;
}

interface Props {
  sortField: string;
  columns: ColumnItem[];
  handleSortingChange: (value: string) => void;
}

const RenderTableHeader = memo(({
  sortField,
  columns,
  handleSortingChange,
}: Props) => {
  return (
    <TableHeader>
      <TableRow>
        {columns.map((item) => (
          <TableCell key={item.accessor}>
            <div className="flex items-center gap-2">
              {item.label}
              {sortField === item.accessor ? (
                <ArrowActive
                  onClick={() => handleSortingChange(item.accessor)}
                />
              ) : (
                <ArrowUnActive
                  onClick={() => handleSortingChange(item.accessor)}
                />
              )}
            </div>
          </TableCell>
        ))}
        <TableCell />
      </TableRow>
    </TableHeader>
  );
});

export { RenderTableHeader };
