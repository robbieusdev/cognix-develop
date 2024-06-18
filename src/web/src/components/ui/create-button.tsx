import { PlusIcon } from 'lucide-react';
import { Button } from './button';

export function CreateButton({ onClick, children, ...props }: React.ButtonHTMLAttributes<HTMLButtonElement>) {
  return (
    <Button variant={'outline'} onClick={onClick} {...props}>
      <PlusIcon className="mr-2 h-4 w-4" />
      {children}
    </Button>
  );
}
