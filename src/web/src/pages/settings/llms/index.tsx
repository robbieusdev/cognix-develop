import { Tabs, TabsContent } from "@/components/ui/tabs";
import { RenderTable } from "@/components/renderTable/render-table";
import { SettingHeader } from "@/components/ui/setting-header";
import { Controller } from "./llm.controller";
import { ConfirmDeleteDialog } from "@/components/dialogs/ConfirmDeleteDialog";
import { useEffect, useState } from "react";
import { LLMDialog } from "@/components/dialogs/LLMDialog";
import axios from "axios";
import { Persona } from "@/models/settings";

export function LLMManagementComponent() {
  const [llms, setLlms] = useState<Persona[]>([]);
  const [selectedRow, setSelectedRow] = useState<Persona>();

  const { columns, sortField, handleSortingChange } =
    Controller.useFilterHandler(llms);

  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showLLMDialogOpen, setShowLLMDialogOpen] = useState(false);

  async function getLLMs() {
    await axios
      .get(import.meta.env.VITE_PLATFORM_API_LLM_LIST_URL)
      .then(function (response) {
        if (response.status == 200) {
          setLlms(response.data.data);
        } else {
          setLlms([]);
        }
      })
      .catch(function (error) {
        console.error("Error fetching messages:", error);
      });
  }

  async function deleteLLM(id: string) {
    await axios
      .post(`${import.meta.env.VITE_PLATFORM_API_LLM_DELETE_URL}/${id}/delete`)
      .then((response) => {
        if (response.status == 200) {
          setLlms(llms.filter((llm) => llm.id !== id));
        }
      });
  }

  async function setRow(id: string) {
    return await axios
      .get(`${import.meta.env.VITE_PLATFORM_API_LLM_DETAIL_URL}/${id}`)
      .then(function (response) {
        if (response.status === 200) {
          setSelectedRow(response.data.data);
        } else {
          throw new Error("Failed to fetch detailed LLM: " + response.status);
        }
      })
      .catch(function (error) {
        console.error("Error fetching detailed LLM:", error);
        throw error; // Re-throw the error to be caught by the caller
      });
  }

  useEffect(() => {
    getLLMs();
  }, [showLLMDialogOpen, showDeleteDialog]);

  return (
    <div className="flex flex-grow flex-col m-8 overflow-x-hidden no-scrollbar">
      <SettingHeader
        title={"LLMs"}
        buttonTitle="New LLM"
        withBtn
        handleClick={() => {
          setShowLLMDialogOpen(true);
        }}
      />
      <>
        <Tabs defaultValue="personal">
          <TabsContent value="personal">
            <RenderTable
              columns={columns}
              handleSortingChange={handleSortingChange}
              sortField={sortField}
              tableData={llms}
              onDelete={async (id: string) => {
                await setRow(id);
                setShowDeleteDialog(true);
              }}
              onEdit={async (id: string) => {
                await setRow(id);
                setShowLLMDialogOpen(true);
              }}
              withBtn
            />
          </TabsContent>
        </Tabs>
      </>
      {showDeleteDialog && (
        <div className="ml-auto">
          <ConfirmDeleteDialog
            description="Are you sure you want to delete this LLM?"
            deleteButtonText="Yes, Delete"
            onConfirm={() => {
              deleteLLM(selectedRow!.id);
            }}
            open={showDeleteDialog}
            onOpenChange={setShowDeleteDialog}
          />
        </div>
      )}
      {showLLMDialogOpen && (
        <LLMDialog
          open={showLLMDialogOpen}
          onOpenChange={() => {
            setShowLLMDialogOpen(false);
            setSelectedRow(undefined);
          }}
          instance={selectedRow}
        />
      )}
    </div>
  );
}

export { LLMManagementComponent as Component };
