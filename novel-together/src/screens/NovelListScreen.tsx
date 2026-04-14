import React, { useState } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, TextInput } from 'react-native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';

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

type NovelListScreenProps = {
  navigation: NativeStackNavigationProp<RootStackParamList, 'NovelList'>;
};

interface Novel {
  id: string;
  title: string;
  author: string;
  description: string;
  chapterCount: number;
}

// Mock data for demonstration
const mockNovels: Novel[] = [
  { id: '1', title: 'The Great Adventure', author: 'John Doe', description: 'An epic journey through magical lands', chapterCount: 15 },
  { id: '2', title: 'Mystery of the Night', author: 'Jane Smith', description: 'A thrilling detective story', chapterCount: 8 },
  { id: '3', title: 'Love in Paris', author: 'Emily Brown', description: 'A romantic tale in the city of love', chapterCount: 20 },
];

export default function NovelListScreen({ navigation }: NovelListScreenProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [novels] = useState<Novel[]>(mockNovels);

  const filteredNovels = novels.filter(novel =>
    novel.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    novel.author.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const renderNovelItem = ({ item }: { item: Novel }) => (
    <TouchableOpacity
      style={styles.novelItem}
      onPress={() => navigation.navigate('NovelDetail', { novelId: item.id })}
    >
      <Text style={styles.novelTitle}>{item.title}</Text>
      <Text style={styles.novelAuthor}>by {item.author}</Text>
      <Text style={styles.novelDescription}>{item.description}</Text>
      <Text style={styles.chapterCount}>{item.chapterCount} chapters</Text>
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Browse Novels</Text>
      
      <TextInput
        style={styles.searchInput}
        placeholder="Search novels..."
        value={searchQuery}
        onChangeText={setSearchQuery}
      />
      
      <FlatList
        data={filteredNovels}
        renderItem={renderNovelItem}
        keyExtractor={item => item.id}
        contentContainerStyle={styles.list}
      />
      
      <TouchableOpacity
        style={styles.backButton}
        onPress={() => navigation.navigate('Home')}
      >
        <Text style={styles.backButtonText}>Back to Home</Text>
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
    marginBottom: 20,
    textAlign: 'center',
  },
  searchInput: {
    borderWidth: 1,
    borderColor: '#ccc',
    borderRadius: 5,
    padding: 10,
    marginBottom: 20,
    fontSize: 16,
  },
  list: {
    paddingBottom: 20,
  },
  novelItem: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 15,
    marginBottom: 15,
    backgroundColor: '#f9f9f9',
  },
  novelTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  novelAuthor: {
    fontSize: 14,
    color: '#666',
    marginBottom: 5,
  },
  novelDescription: {
    fontSize: 14,
    color: '#333',
    marginBottom: 5,
  },
  chapterCount: {
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
