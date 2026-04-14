import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, TextInput, ActivityIndicator, Alert } from 'react-native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { novelApi } from '../services/api';

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
  description: string;
  author_username: string;
  chapter_count: number;
}

export default function NovelListScreen({ navigation }: NovelListScreenProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [novels, setNovels] = useState<Novel[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchNovels = async () => {
    setLoading(true);
    try {
      const response = await novelApi.list();
      setNovels(response.data);
    } catch (error: any) {
      Alert.alert('Error', 'Failed to load novels');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNovels();
  }, []);

  const filteredNovels = novels.filter(novel =>
    novel.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    novel.author_username?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const renderNovelItem = ({ item }: { item: Novel }) => (
    <TouchableOpacity
      style={styles.novelItem}
      onPress={() => navigation.navigate('NovelDetail', { novelId: item.id })}
    >
      <Text style={styles.novelTitle}>{item.title}</Text>
      <Text style={styles.novelAuthor}>by {item.author_username}</Text>
      <Text style={styles.novelDescription}>{item.description}</Text>
      <Text style={styles.chapterCount}>{item.chapter_count} chapters</Text>
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      <Text style={styles.title}>浏览小说</Text>
      
      <TextInput
        style={styles.searchInput}
        placeholder="搜索..."
        value={searchQuery}
        onChangeText={setSearchQuery}
      />
      
      {loading ? (
        <ActivityIndicator size="large" />
      ) : (
        <FlatList
          data={filteredNovels}
          renderItem={renderNovelItem}
          keyExtractor={item => item.id}
          contentContainerStyle={styles.list}
          ListEmptyComponent={<Text style={styles.empty}>暂无小说</Text>}
        />
      )}
      
      <TouchableOpacity
        style={styles.backButton}
        onPress={() => navigation.navigate('Home')}
      >
        <Text style={styles.backButtonText}>返回首页</Text>
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
  empty: {
    textAlign: 'center',
    color: '#666',
    marginTop: 20,
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