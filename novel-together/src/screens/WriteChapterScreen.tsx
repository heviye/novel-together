import React, { useState } from 'react';
import { View, Text, TextInput, Button, StyleSheet, Alert, ScrollView } from 'react-native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { RouteProp } from '@react-navigation/native';

type RootStackParamList = {
  Home: undefined;
  Login: undefined;
  Register: undefined;
  NovelList: undefined;
  NovelDetail: { novelId: string };
  WriteChapter: { novelId: string };
  Chapter: { chapterId: string };
  Profile: undefined;
};

type WriteChapterScreenProps = {
  navigation: NativeStackNavigationProp<RootStackParamList, 'WriteChapter'>;
  route: RouteProp<RootStackParamList, 'WriteChapter'>;
};

export default function WriteChapterScreen({ navigation, route }: WriteChapterScreenProps) {
  const { novelId } = route.params;
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');

  const handleSave = () => {
    if (!title || !content) {
      Alert.alert('Error', 'Please fill in all fields');
      return;
    }
    // TODO: Implement actual save chapter API call
    Alert.alert('Success', 'Chapter saved successfully');
    navigation.goBack();
  };

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.title}>Write New Chapter</Text>
      <Text style={styles.novelId}>Novel ID: {novelId}</Text>
      
      <View style={styles.form}>
        <Text style={styles.label}>Chapter Title</Text>
        <TextInput
          style={styles.input}
          placeholder="Enter chapter title"
          value={title}
          onChangeText={setTitle}
        />
        
        <Text style={styles.label}>Content</Text>
        <TextInput
          style={[styles.input, styles.contentInput]}
          placeholder="Write your story here..."
          value={content}
          onChangeText={setContent}
          multiline
          textAlignVertical="top"
        />
        
        <View style={styles.buttonContainer}>
          <Button title="Save Chapter" onPress={handleSave} />
          <Button title="Cancel" onPress={() => navigation.goBack()} />
        </View>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  novelId: {
    fontSize: 14,
    color: '#666',
    marginBottom: 20,
  },
  form: {
    gap: 15,
  },
  label: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  input: {
    borderWidth: 1,
    borderColor: '#ccc',
    borderRadius: 5,
    padding: 10,
    fontSize: 16,
  },
  contentInput: {
    minHeight: 300,
  },
  buttonContainer: {
    gap: 10,
    marginTop: 20,
  },
});
