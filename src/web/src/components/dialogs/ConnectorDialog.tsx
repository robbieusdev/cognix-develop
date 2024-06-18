import { Button, buttonVariants } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { DefaultValues, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input";
import { useMutation } from "@/lib/mutation";
import { toast } from "react-toastify";
import {
  CreateConnectorSchema,
  UpdateConnectorSchema,
} from "@/lib/schemas/connectors";
import { TextArea } from "../ui/textarea";
import { Connector, SourceType } from "@/models/settings";
import { useEffect, useState } from "react";
import axios from "axios";
import { capitalize } from "@/lib/utils";
import { formSchema } from "@/lib/validations/connectors";

export function CreateConnectorDialog({
  defaultValues,
  instance,
  children,
  open,
  onOpenChange,
}: {
  defaultValues?: DefaultValues<z.infer<typeof formSchema>>;
  instance?: Connector;
  children?: React.ReactNode;
  open?: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const [step, setStep] = useState(1);
  const [sourceTypes, setSourceTypes] = useState<SourceType[]>([]);
  const [loading, setLoading] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: instance?.name || "",
      connector_specific_config: instance
        ? JSON.stringify(
            instance.connector_specific_config ?? '{"your_key": "your_value"}'
          )
        : "",
      refresh_freq: instance ? String(instance.refresh_freq) : "",
      source: instance?.source || "",
      credential_id: instance?.credential_id || "",
      ...defaultValues,
    },
  });

  const { trigger: triggerCreateConnector } =
    useMutation<CreateConnectorSchema>(
      import.meta.env.VITE_PLATFORM_API_CONNECTOR_CREATE_URL,
      "POST"
    );

  const { trigger: triggerEditConnector } = useMutation<UpdateConnectorSchema>(
    `${import.meta.env.VITE_PLATFORM_API_CONNECTOR_EDIT_URL}/${instance?.id}`,
    "PUT"
  );

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    const data = {
      connector_specific_config: JSON.parse(values.connector_specific_config),
      refresh_freq: Number(values.refresh_freq),
      name: values.name,
      source: values.source,
      credential_id: values.credential_id ?? "0",
    };
    if (instance) {
      try {
        await triggerEditConnector(data);
        toast.success("Connector successfully updated");
      } catch (e) {
        toast.error(e as string);
      }
    } else {
      try {
        await triggerCreateConnector(data);
        toast.success("Connector successfully created");
      } catch (e) {
        toast.error(e as string);
      }
    }
    onOpenChange(false);
  };

  async function getSourceTypes() {
    setLoading(true);
    await axios
      .get(import.meta.env.VITE_PLATFORM_API_SOURCE_TYPES_LIST_URL)
      .then((response) => {
        if (response.status == 200) {
          setSourceTypes(response.data.data);
        }
      })
      .finally(() => {
        setLoading(false);
      });
    return [];
  }

  useEffect(() => {
    if (step == 1 && !instance) {
      getSourceTypes();
    }
  }, []);

  return (
    <Dialog
      open={open}
      onOpenChange={(newOpen) => {
        onOpenChange?.(newOpen);
        form.reset();
      }}
    >
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent
        className={step == 1 ? "sm:max-w-[800px]" : "sm:max-w-[425px]"}
      >
        <DialogHeader>
          <DialogTitle>
            {instance ? "Edit Connector" : "Add Connector"}
          </DialogTitle>
          {!instance && (
            <span className="fixed right-3 top-3 text-sm">
              <span className="text-primary">Step {step}</span>/2
            </span>
          )}
        </DialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="max-w-full space-y-4 overflow-hidden px-0.5"
            onKeyDown={(e) => e.key === "Enter" && e.preventDefault()}
          >
            {step === 1 && !instance ? (
              <div className="sm:max-w-[800px] space-y-5">
                <FormField
                  control={form.control}
                  name="source"
                  render={() => (
                    <>
                      <FormItem>
                        <FormControl>
                          {loading ? (
                            <div className="fixed top-0 left-0 w-full h-full flex items-center justify-center bg-gray-200 bg-opacity-50 z-50">
                              <div className="animate-spin rounded-full h-10 w-10 border-t-4 border-b-4"></div>
                            </div>
                          ) : (
                            <div className="grid grid-cols-4 gap-4 border-none">
                              {sourceTypes.map((sourceType) => (
                                <div
                                  key={sourceType.id}
                                  onClick={() =>
                                    form.setValue("source", sourceType.id)
                                  }
                                  className={`bg-white p-4 rounded-sm border ${
                                    form.getValues("source") === sourceType.id
                                      ? "border-primary"
                                      : ""
                                  } cursor-pointer`}
                                >
                                  <p className="text-gray-600 text-center">
                                    {sourceType.name}
                                  </p>
                                </div>
                              ))}
                            </div>
                          )}
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                      <DialogFooter>
                        <DialogClose asChild>
                          <Button
                            variant="outline"
                            className={`${buttonVariants({
                              variant: "destructive",
                            })} w-full h-10`}
                          >
                            Cancel
                          </Button>
                        </DialogClose>
                        <Button
                          onClick={() => {
                            if (form.getValues("source")) {
                              setStep(2);
                            }
                          }}
                          className="w-full h-10"
                        >
                          Next
                        </Button>
                      </DialogFooter>
                    </>
                  )}
                />
              </div>
            ) : (
              <>
                <FormField
                  control={form.control}
                  name="source"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Connector</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          value={capitalize(field.value)}
                          disabled
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => {
                    return (
                      <FormItem>
                        <FormControl>
                          <Input placeholder="Name" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    );
                  }}
                />

                <FormField
                  control={form.control}
                  name="connector_specific_config"
                  render={({ field }) => (
                    <FormItem>
                      <FormControl>
                        <TextArea
                          placeholder={`Connector Specific Configuration \n {"your_key": "your_value"}`}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="refresh_freq"
                  render={({ field }) => (
                    <FormItem>
                      <FormControl>
                        <Input
                          placeholder="Refresh Frequency (seconds)"
                          type="number"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                {!instance && (
                  <FormField
                    control={form.control}
                    name="credential_id"
                    render={({ field }) => (
                      <FormItem>
                        <FormControl>
                          <TextArea
                            placeholder="Connector credential"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                )}
                <DialogFooter>
                  <DialogClose asChild>
                    <Button
                      variant="outline"
                      className={`${buttonVariants({
                        variant: "destructive",
                      })} w-full h-10`}
                    >
                      Cancel
                    </Button>
                  </DialogClose>
                  <Button
                    type="submit"
                    disabled={!form.formState.errors}
                    className="w-full h-10"
                  >
                    {instance ? "Update" : "Add"}
                  </Button>
                </DialogFooter>
              </>
            )}
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
