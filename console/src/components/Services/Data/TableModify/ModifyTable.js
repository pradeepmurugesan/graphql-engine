import PropTypes from 'prop-types';
import React from 'react';
import TableHeader from '../TableCommon/TableHeader';

import { getAllDataTypeMap } from '../Common/utils';

import { TABLE_ENUMS_SUPPORT } from '../../../../helpers/versionUtils';
import globals from '../../../../Globals';

import {
  deleteTableSql,
  untrackTableSql,
  RESET,
  setUniqueKeys,
  toggleTableAsEnum,
} from '../TableModify/ModifyActions';
import {
  setTable,
  fetchColumnTypeInfo,
  RESET_COLUMN_TYPE_INFO,
} from '../DataActions';
import Button from '../../../Common/Button/Button';
import ColumnEditorList from './ColumnEditorList';
import ColumnCreator from './ColumnCreator';
import PrimaryKeyEditor from './PrimaryKeyEditor';
import TableCommentEditor from './TableCommentEditor';
import EnumsSection, {
  EnumTableModifyWarning,
} from '../Common/ReusableComponents/EnumsSection';
import ForeignKeyEditor from './ForeignKeyEditor';
import UniqueKeyEditor from './UniqueKeyEditor';
import TriggerEditorList from './TriggerEditorList';
import styles from './ModifyTable.scss';
import { NotFoundError } from '../../../Error/PageNotFound';

import { getConfirmation } from '../../../Common/utils/jsUtils';

class ModifyTable extends React.Component {
  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ type: RESET });
    dispatch(setTable(this.props.tableName));
    dispatch(fetchColumnTypeInfo());
  }
  componentWillUnmount() {
    this.props.dispatch({
      type: RESET_COLUMN_TYPE_INFO,
    });
  }
  render() {
    const {
      tableName,
      allSchemas,
      dispatch,
      migrationMode,
      currentSchema,
      tableCommentEdit,
      columnEdit,
      pkModify,
      fkModify,
      dataTypes,
      validTypeCasts,
      uniqueKeyModify,
      columnDefaultFunctions,
      schemaList,
      tableEnum,
    } = this.props;

    const dataTypeIndexMap = getAllDataTypeMap(dataTypes);

    const tableSchema = allSchemas.find(
      t => t.table_name === tableName && t.table_schema === currentSchema
    );
    if (!tableSchema) {
      // throw a 404 exception
      throw new NotFoundError();
    }
    const tableComment = tableSchema.comment;

    const untrackBtn = (
      <Button
        type="submit"
        className={styles.add_mar_right}
        color="white"
        size="sm"
        onClick={() => {
          const confirmMessage = `This will remove the table "${tableName}" from the GraphQL schema`;
          const isOk = getConfirmation(confirmMessage);
          if (isOk) {
            dispatch(untrackTableSql(tableName));
          }
        }}
        data-test="untrack-table"
      >
        Untrack Table
      </Button>
    );

    const deleteBtn = (
      <Button
        type="submit"
        color="red"
        size="sm"
        onClick={() => {
          const confirmMessage = `This will permanently delete the table "${tableName}" from the database`;
          const isOk = getConfirmation(confirmMessage, true, tableName);
          if (isOk) {
            dispatch(deleteTableSql(tableName, tableSchema));
          }
        }}
        data-test="delete-table"
      >
        Delete table
      </Button>
    );

    const getEnumsSection = () => {
      const supportEnums =
        globals.featuresCompatibility &&
        globals.featuresCompatibility[TABLE_ENUMS_SUPPORT];
      if (!supportEnums) return null;

      const toggleEnum = () => dispatch(toggleTableAsEnum(tableSchema.is_enum));

      return (
        <React.Fragment>
          <EnumsSection
            isEnum={tableSchema.is_enum}
            toggleEnum={toggleEnum}
            loading={tableEnum.loading}
          />
          <hr />
        </React.Fragment>
      );
    };

    // if (tableSchema.primary_key.columns > 0) {}
    return (
      <div className={`${styles.container} container-fluid`}>
        <TableHeader
          dispatch={dispatch}
          table={tableSchema}
          tabName="modify"
          migrationMode={migrationMode}
        />
        <br />
        <div className={`container-fluid ${styles.padd_left_remove}`}>
          <div
            className={
              `col-xs-10 ${styles.padd_left_remove}` +
              ' ' +
              styles.modifyMinWidth
            }
          >
            <TableCommentEditor
              tableComment={tableComment}
              tableCommentEdit={tableCommentEdit}
              isTable
              dispatch={dispatch}
            />
            <EnumTableModifyWarning isEnum={tableSchema.is_enum} />
            <h4 className={styles.subheading_text}>Columns</h4>
            <ColumnEditorList
              validTypeCasts={validTypeCasts}
              dataTypeIndexMap={dataTypeIndexMap}
              tableSchema={tableSchema}
              columnEdit={columnEdit}
              dispatch={dispatch}
              currentSchema={currentSchema}
              columnDefaultFunctions={columnDefaultFunctions}
            />
            <hr />
            <h4 className={styles.subheading_text}>Add a new column</h4>
            <ColumnCreator
              dispatch={dispatch}
              tableName={tableName}
              dataTypes={dataTypes}
              validTypeCasts={validTypeCasts}
              columnDefaultFunctions={columnDefaultFunctions}
            />
            <hr />
            <h4 className={styles.subheading_text}>Primary Key</h4>
            <PrimaryKeyEditor
              tableSchema={tableSchema}
              pkModify={pkModify}
              dispatch={dispatch}
              currentSchema={currentSchema}
            />
            <hr />
            <h4 className={styles.subheading_text}>Foreign Keys</h4>
            <ForeignKeyEditor
              tableSchema={tableSchema}
              currentSchema={currentSchema}
              allSchemas={allSchemas}
              schemaList={schemaList}
              dispatch={dispatch}
              fkModify={fkModify}
            />
            <hr />
            <h4 className={styles.subheading_text}>Unique Keys</h4>
            <UniqueKeyEditor
              tableSchema={tableSchema}
              currentSchema={currentSchema}
              allSchemas={allSchemas}
              dispatch={dispatch}
              uniqueKeys={uniqueKeyModify}
              setUniqueKeys={setUniqueKeys}
            />
            <hr />
            <h4 className={styles.subheading_text}>Triggers</h4>
            <TriggerEditorList tableSchema={tableSchema} dispatch={dispatch} />
            <hr />
            {getEnumsSection()}
            {untrackBtn}
            {deleteBtn}
            <br />
            <br />
          </div>
        </div>
      </div>
    );
  }
}

ModifyTable.propTypes = {
  tableName: PropTypes.string.isRequired,
  currentSchema: PropTypes.string.isRequired,
  allSchemas: PropTypes.array.isRequired,
  migrationMode: PropTypes.bool.isRequired,
  activeEdit: PropTypes.object.isRequired,
  fkAdd: PropTypes.object.isRequired,
  relAdd: PropTypes.object.isRequired,
  ongoingRequest: PropTypes.bool.isRequired,
  lastError: PropTypes.object,
  lastFormError: PropTypes.object,
  columnEdit: PropTypes.object.isRequired,
  lastSuccess: PropTypes.bool,
  dispatch: PropTypes.func.isRequired,
  pkModify: PropTypes.array.isRequired,
  fkModify: PropTypes.array.isRequired,
  serverVersion: PropTypes.string,
};

const mapStateToProps = (state, ownProps) => ({
  tableName: ownProps.params.table,
  allSchemas: state.tables.allSchemas,
  migrationMode: state.main.migrationMode,
  serverVersion: state.main.serverVersion,
  currentSchema: state.tables.currentSchema,
  columnEdit: state.tables.modify.columnEdit,
  pkModify: state.tables.modify.pkModify,
  fkModify: state.tables.modify.fkModify,
  dataTypes: state.tables.columnDataTypes,
  columnDefaultFunctions: state.tables.columnDefaultFunctions,
  validTypeCasts: state.tables.columnTypeCasts,
  columnDataTypeFetchErr: state.tables.columnDataTypeFetchErr,
  schemaList: state.tables.schemaList,
  ...state.tables.modify,
});

const modifyTableConnector = connect => connect(mapStateToProps)(ModifyTable);

export default modifyTableConnector;
