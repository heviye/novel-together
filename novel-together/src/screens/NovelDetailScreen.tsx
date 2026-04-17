import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, ActivityIndicator, Alert, Button } from 'react-native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { RouteProp } from '@react-navigation/native';
import { novelApi, chapterApi } from '../services/api';

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

type Novel = {
  id: string;
  title: string;
  author: string;
  description: string;
  genre: string;
  status: string;
};

type Chapter = {
  id: string;
  title: string;
  number: number;
  author: string;
};

// Mock data fallback
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
  const [novel, setNovel] = useState<Novel | null>(null);
  const [chapters, setChapters] = useState<Chapter[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchNovelData();
  }, [novelId]);

  const fetchNovelData = async () => {
    setLoading(true);
    setError(null);
    try {
      // Fetch novel details
      const novelResponse = await novelApi.get(novelId);
      setNovel(novelResponse.data);
      
      // Fetch chapters
      const chaptersResponse = await novelApi.getChapters(novelId);
      setChapters(chaptersResponse.data);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load novel data');
      Alert.alert('Error', 'Failed to load novel data, using mock data');
      // Fallback to mock data
      setNovel(mockNovel);
      setChapters(mockChapters);
    } finally {
      setLoading(false);
    }
  };

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

  if (loading) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  if (!novel) {
    return null;
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>{novel.title}</Text>
      <Text style={styles.author}>by {novel.author}</Text>
      <View style={styles.metaContainer}>
        <Text style={styles.genre}>{novel.genre}</Text>
        <Text style={styles.status}>{novel.status}</Text>
      </View>
      <Text style={styles.description}>{novel.description}</Text>
      
      <Button
        title="Write New Chapter"
        onPress={() => navigation.navigate('WriteChapter', { novelId: novel.id })}
      />
      
      <Text style={styles.chaptersTitle}>Chapters</Text>
      <FlatList
        data={chapters}
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
    marginBottom: 10,
  },
  author: {
    fontSize: 18,
    color: '#333',
    marginBottom: 5,
  },
  metaContainer: {
    flexDirection: 'row',
    marginBottom: 15,
  },
  genre: {
    fontSize: 14,
    color: '#666',
    marginRight: 15,
  },
  status: {
    fontSize: 14,
    color: '#666',
  },
  description: {
    fontSize: 16,
    color: '#333',
    marginBottom: 20,
    lineHeight: 24,
  },
  chaptersTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 15,
  },
  chapterItem: {
    padding: 15,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  chapterNumber: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  chapterTitle: {
    fontSize: 16,
    marginTop: 5,
  },
  chapterAuthor: {
    fontSize: 14,
    color: '#666',
    marginTop: 5,
  },
  backButton: {
    padding: 15,
    backgroundColor: '#f0f0f0',
    borderRadius: 5,
    alignItems: 'center',
    marginTop: 10,
  },
  homeButton: {
    padding: 15,
  },
  backButtonText: {
    color: '#007AFF',
    fontWeight: 'bold',
  },
  homeButtonText: {
    color: '#007AFF',
    fontWeight: 'bold',
  },
});
