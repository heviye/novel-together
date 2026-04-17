import React from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, Button } from 'react-native';
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

type ChapterScreenProps = {
  navigation: NativeStackNavigationProp<RootStackParamList, 'Chapter'>;
  route: RouteProp<RootStackParamList, 'Chapter'>;
};

// Mock data for demonstration
const mockChapter = {
  id: '1',
  title: 'The Beginning',
  number: 1,
  author: 'John Doe',
  content: `It was a dark and stormy night when the young protagonist first set out on their journey. The village behind them was silent, its inhabitants sleeping soundly, unaware of the adventure that awaited.

The path ahead was shrouded in mist, but our hero pressed forward with determination. Every step brought them closer to their destiny, closer to uncovering the ancient secrets that had been hidden for centuries.

As they traveled through the forest, strange sounds echoed around them. Twigs snapped underfoot, and shadows seemed to dance between the trees. But our hero was not afraid. They had trained for this moment their entire life.

Finally, they reached the edge of the forest. Before them lay a vast landscape of mountains and valleys, each promising new challenges and new discoveries. With a deep breath, our hero stepped forward into the unknown.`,
  createdAt: '2024-01-15',
};

export default function ChapterScreen({ navigation, route }: ChapterScreenProps) {
  const { chapterId } = route.params;

  return (
    <View style={styles.container}>
      <ScrollView style={styles.scrollView}>
        <Text style={styles.title}>{mockChapter.title}</Text>
        <Text style={styles.meta}>
          Chapter {mockChapter.number} | by {mockChapter.author}
        </Text>
        <Text style={styles.date}>Published: {mockChapter.createdAt}</Text>
        
        <Text style={styles.content}>{mockChapter.content}</Text>
      </ScrollView>
      
      <View style={styles.buttonContainer}>
        <Button
          title="Back to Novel"
          onPress={() => navigation.goBack()}
        />
        <Button
          title="Home"
          onPress={() => navigation.navigate('Home')}
        />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  scrollView: {
    flex: 1,
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  meta: {
    fontSize: 14,
    color: '#666',
    marginBottom: 5,
  },
  date: {
    fontSize: 12,
    color: '#999',
    marginBottom: 20,
  },
  content: {
    fontSize: 16,
    lineHeight: 28,
    color: '#333',
  },
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    padding: 15,
    borderTopWidth: 1,
    borderTopColor: '#ddd',
  },
});
