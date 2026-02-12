import { colors } from "@/constants";
import { ForwardedRef, forwardRef, ReactNode } from "react";
import {
  StyleSheet,
  Text,
  TextInput,
  TextInputProps,
  View,
} from "react-native";
import hairlineWidth = StyleSheet.hairlineWidth;

interface InputFieldProps extends TextInputProps {
  label?: string;
  variant?: "filled" | "standard" | "outlined";
  error?: string;
  customHeight?: number;
  rightChild?: ReactNode;
}

function InputField(
  {
    label,
    variant = "filled",
    error = "",
    rightChild = null,
    customHeight = 44,
    ...props
  }: InputFieldProps,
  ref?: ForwardedRef<TextInput>,
) {
  return (
    <View>
      {label && <Text style={styles.label}>{label}</Text>}
      <View
        style={[
          styles.container,
          styles[variant],
          props.multiline && styles.multiline,
          Boolean(error) && styles.inputError,
          { height: customHeight },
        ]}
      >
        <TextInput
          ref={ref}
          autoCapitalize="none"
          placeholderTextColor={colors.GRAY_400}
          spellCheck={false}
          autoCorrect={false}
          {...props}
          style={[styles.input, styles[`${variant}Text`], props.style]}
        />
        {rightChild}
      </View>
      {Boolean(error) && <Text style={styles.error}>{error}</Text>}
    </View>
  );
}

const styles = StyleSheet.create({
  label: {
    fontSize: 12,
    color: colors.GRAY_700,
    marginBottom: 5,
  },
  container: {
    height: 44,
    borderRadius: 8,
    paddingHorizontal: 10,
    alignItems: "center",
    justifyContent: "center",
    flexDirection: "row",
  },
  filled: {
    backgroundColor: colors.GRAY_100,
  },
  standard: {
    borderWidth: 1,
    borderColor: colors.GRAY_200,
    backgroundColor: colors.WHITE,
  },
  outlined: {
    backgroundColor: colors.WHITE,
    borderWidth: hairlineWidth,
    borderColor: colors.GRAY_700,
  },
  standardText: {
    color: colors.BLACK,
  },
  outlinedText: {
    color: colors.BLACK,
  },
  filledText: {},
  input: {
    fontSize: 16,
    padding: 0,
    flex: 1,
  },
  error: {
    fontSize: 12,
    marginTop: 5,
    color: colors.RED_500,
  },
  inputError: {
    backgroundColor: colors.RED_100,
  },
  multiline: {
    alignItems: "flex-start",
    paddingVertical: 10,
    height: 200,
  },
});

export default forwardRef(InputField);
