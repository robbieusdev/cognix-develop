export enum LABEL_STATUS {
  STATUS = 'Status',
  SUCCEEDED = 'Succeeded',
  ERROR = 'Error',
  ACTIVE = 'Active',
  DISABLED = 'Disable',
  VISIBLE = 'Visible',
  USER = 'User',
  ADMIN = 'Admin',
  SUPER_ADMIN = 'Super Admin',
}

export interface Props {
  label: LABEL_STATUS | string | number;
}

const LabelStatus = ({ label }: Props) => {
  const labelContainerStyles: Record<Props['label'], string> = {
    [LABEL_STATUS.ACTIVE]: '#E6F7F2',
    [LABEL_STATUS.ADMIN]: '#FDE7DB',
    [LABEL_STATUS.DISABLED]: '#F4F4F4',
    [LABEL_STATUS.ERROR]: '#FDE7DB',

    [LABEL_STATUS.STATUS]: '#FDE7DB',
    [LABEL_STATUS.SUCCEEDED]: '#E6F7F2',
    [LABEL_STATUS.SUPER_ADMIN]: '#F2E8F7',
    [LABEL_STATUS.VISIBLE]: '#E6F7F2',
    [LABEL_STATUS.USER]: '#DAF4FF',
  };

  const labelTextStyles: Record<Props['label'], string> = {
    [LABEL_STATUS.ACTIVE]: '#047C57',
    [LABEL_STATUS.STATUS]: '#F0681C',
    [LABEL_STATUS.SUCCEEDED]: '#047C57',
    [LABEL_STATUS.ERROR]: '#DB3A34',
    [LABEL_STATUS.DISABLED]: '#9299A3',
    [LABEL_STATUS.ADMIN]: '#F0681C',
    [LABEL_STATUS.VISIBLE]: '#047C57',
    [LABEL_STATUS.SUPER_ADMIN]: '#9C50C7',
    [LABEL_STATUS.USER]: '#006A95',
  };

  return (
    <div
      style={{
        backgroundColor: !labelContainerStyles[label]
          ? 'white'
          : labelContainerStyles[label],
      }}
      className={`text-white gap-3 focus:ring-4 focus:outline-none font-medium rounded-lg text-sm px-2.5 py-1.5 text-center inline-flex items-center`}
    >
      <p
        style={{
          color: !labelTextStyles[label] ? '#272727' : labelTextStyles[label],
        }}
      >
        {label}
      </p>
    </div>
  );
};
LabelStatus.displayName = 'Textarea';

export { LabelStatus };
