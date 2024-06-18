import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTrigger,
  } from '@/components/ui/alert-dialog';
  import { buttonVariants } from '@/components/ui/button';
  
  interface ConfirmDeleteDialogProps {
    description: string;
    deleteButtonText?: string;
    disabled?: boolean;
    onConfirm: () => void;
    open?: boolean;
    onOpenChange?: (open: boolean) => void;
    showTrigger?: boolean;
    children?: React.ReactNode;
  }
  
  export function ConfirmDeleteDialog({
    description,
    deleteButtonText,
    onConfirm,
    showTrigger = true,
    open,
    onOpenChange,
    children,
  }: ConfirmDeleteDialogProps) {
    return (
      <AlertDialog open={open} onOpenChange={onOpenChange}>
        {showTrigger && <AlertDialogTrigger asChild>{children}</AlertDialogTrigger>}
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogDescription>{description}</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel className={`${buttonVariants({ variant: 'destructive' })} w-full`}>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={onConfirm} className='w-full'>
              {deleteButtonText}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    );
  }
  