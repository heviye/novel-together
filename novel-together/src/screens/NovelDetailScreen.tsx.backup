import React from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, Button } from 'react-native';
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

type NovelDetailScreenProps = {
  navigation: NativeStackNavigationProp<RootStackParamList, 'NovelDetail'>;
  route: RouteProp<RootStackParamList, 'NovelDetail'>;
};

interface Chapter {
  id: string;
  title: string;
  number: number;
  author: string;
}

// Mock data for demonstration
const mockNovel = {
  id: '1',
  title: 'The Great Adventure',
  author: 'John Doe',
  description: 'An epic journey through magical lands where heroes rise and fall.',
  genre: 'Fantasy',
  status: 'Ongoing',
};

const mockChapters: Chapter[] = [
  { id: '1', title: 'The Beginning', number: 1, author: 'John Doe' },
  { id: '2', title: 'The First Battle', number: 2, author: 'Jane Smith' },
  { id: '3', title: 'Discovery', number: 3, author: 'John Doe' },
  { id: '4', title: 'New Friends', number: 4, author: 'Jane Smith' },
  { id: '5', title: 'The Challenge', number: 5, author: 'John Doe' },
];

export default function NovelDetailScreen({ navigation, route }: NovelDetailScreenProps) {
  const { novelId } = route.params;

  const renderChapterItem = ({ item }: { item: Chapter }) => (
    <TouchableOpacity
      style={styles.chapterItem}
      onPress={() => navigation.navigate('Chapter', { chapterId: item.id })}
    >
      <Text style={styles.chapterNumber}>Chapter {item.number}</Text>
      <Text style={styles.chapterTitle}>{item.title}</Text>
      <Text style={styles.chapterAuthor}>by {item.author}</Text>
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      <Text style={styles.title}>{mockNovel.title}</Text>
      <Text style={styles.author}>by {mockNovel.author}</Text>
      <View style={styles.metaContainer}>
        <Text style={styles.genre}>{mockNovel.genre}</Text>
        <Text style={styles.status}>{mockNovel.status}</Text>
      </View>
      <Text style={styles.description}>{mockNovel.description}</Text>
      
      <Button
        title="Write New Chapter"
        onPress={() => navigation.navigate('WriteChapter', { novelId })}
      />
      
      <Text style={styles.chaptersTitle}>Chapters</Text>
      <FlatList
        data={mockChapters}
        renderItem={renderChapterItem}
        keyExtractor={item => item.id}
        style={styles.chapterList}
      />
      
      <TouchableOpacity
        style={styles.backButton}
        onPress={() => navigation.navigate('NovelList')}
      >
        <Text style={styles.backButtonText}>Back to Novel List</Text>
      </TouchableOpacity>
    </View>
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
  author: {
    fontSize: 16,
    color: '#666',
    marginBottom: 10,
  },
  metaContainer: {
    flexDirection: 'row',
    gap: 10,
    marginBottom: 15,
  },
  genre: {
    backgroundColor: '#e0e0e0',
    paddingHorizontal: 10,
    paddingVertical: 5,
    borderRadius: 5,
    fontSize: 12,
  },
  status: {
    backgroundColor: '#4CAF50',
    color: '#fff',
    paddingHorizontal: 10,
    paddingVertical: 5,
    borderRadius: 5,
    fontSize: 12,
  },
  description: {
    fontSize: 14,
    color: '#333',
    marginBottom: 20,
    lineHeight: 20,
  },
  chaptersTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginTop: 20,
    marginBottom: 10,
  },
  chapterList: {
    flex: 1,
  },
  chapterItem: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 5,
    padding: 12,
    marginBottom: 10,
    backgroundColor: '#f9f9f9',
  },
  chapterNumber: {
    fontSize: 12,
    color: '#666',
  },
  chapterTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginVertical: 3,
  },
  chapterAuthor: {
    fontSize: 12,
    color: '#007AFF',
  },
  backButton: {
    padding: 15,
    alignItems: 'center',
  },
  backButtonText: {
    color: '#007AFF',
    fontSize: 16,
  },
});
