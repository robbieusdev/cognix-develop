export type CreateConnectorSchema = {
  name: string;
  source: string;
  connector_specific_config: object;
  refresh_freq: number;
  credential_id?: string;
};

export type UpdateConnectorSchema = {
  name: string;
  connector_specific_config: object;
  refresh_freq: number;
  credential_id?: string;
};

export type DisableConnectorSchema = {
  disabled: boolean;
};
