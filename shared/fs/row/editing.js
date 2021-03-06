// @flow
import * as React from 'react'
import * as Types from '../../constants/types/fs'
import {globalStyles, globalColors, globalMargins, platformStyles, isMobile, glamorous} from '../../styles'
import rowStyles from './styles'
import {Icon, ClickableBox, Input, Box, Button, Text} from '../../common-adapters'
import PathItemIcon from '../common/path-item-icon'

type EditingProps = {
  name: string,
  hint: string,
  status: Types.EditStatusType,
  itemStyles: Types.ItemStyles,
  isCreate: boolean,
  onSubmit: () => void,
  onUpdate: (name: string) => void,
  onCancel: () => void,
}

const HoverClickableBox = isMobile
  ? ClickableBox
  : glamorous(ClickableBox)({
      '& .fs-path-item-editing-trash-icon': {
        color: globalColors.black_40,
      },
      '& .fs-path-item-editing-trash-icon:hover': {
        color: globalColors.black_60,
      },
    })

const Editing = (props: EditingProps) => (
  <Box style={rowStyles.rowBox}>
    <PathItemIcon spec={props.itemStyles.iconSpec} style={rowStyles.pathItemIcon} />
    <Box key="main" style={rowStyles.itemBox}>
      <Input
        hideUnderline={true}
        small={true}
        value={props.name}
        hintText={props.hint}
        inputStyle={stylesText}
        onEnterKeyDown={props.onSubmit}
        onChangeText={name => props.onUpdate(name)}
        autoFocus={true}
        selectTextOnFocus={true}
      />
    </Box>
    <Box key="right" style={rowStyles.rightBox}>
      {props.status === 'failed' && <Text type="BodyError">Failed</Text>}
      <Button
        key="create"
        style={stylesButton}
        type="Primary"
        small={true}
        label={props.status === 'failed' ? 'Retry' : props.isCreate ? 'Create' : 'Save'}
        waiting={props.status === 'saving'}
        onClick={props.status === 'saving' ? undefined : props.onSubmit}
      />
      <HoverClickableBox style={stylesCancelBox} onClick={props.onCancel}>
        <Icon type="iconfont-trash" className="fs-path-item-editing-trash-icon" style={stylesIconCancel} />
      </HoverClickableBox>
    </Box>
  </Box>
)

const stylesCancelBox = {
  ...globalStyles.flexBoxRow,
  alignItems: 'center',
}

const stylesIconCancel = platformStyles({
  common: {
    padding: globalMargins.tiny,
    paddingRight: 0,
  },
  isMobile: {
    fontSize: 22,
  },
})

const stylesText = platformStyles({
  common: {
    ...globalStyles.fontSemibold,
    maxWidth: '100%',
  },
  isMobile: {
    marginTop: 22,
  },
})

const stylesButton = {
  marginLeft: globalMargins.tiny,
}

export default Editing
